package ec2

import (
	//"encoding/json"
	//"io/ioutil"
	//"net/http"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	//"github.com/davecgh/go-spew/spew"
)

var client *ec2.EC2

func Setup() {
	client = ec2.New(session.New())
}

func DescribeInstances() ([]*ec2.Instance, error) {
	var filters []*ec2.Filter

	filters = autofilters()

	params := &ec2.DescribeInstancesInput{Filters: filters}
	instances := make([]*ec2.Instance, 0)

	err := client.DescribeInstancesPages(params,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, res := range page.Reservations {
				instances = append(instances, res.Instances...)
			}
			return true
		})

	return instances, err
}

func getTagValue(inst *ec2.Instance, tagKey string) *string {
	for _, tag := range inst.Tags {
		if *tag.Key == tagKey {
			return tag.Value
		}
	}

	return nil
}

func instanceName(inst *ec2.Instance) *string {
	ip := inst.PrivateIpAddress
	//spew.Dump(ip)

	if ip == nil {
		return nil
	}

	ip2 := strings.Replace(*ip, ".", "-", -1)
	ip2ptr := fmt.Sprintf("%s-ip-%s", *getTagValue(inst, "Name"), ip2)

	return &ip2ptr
}

func InstanceNames() (map[string]bool, error) {
	instances, err := DescribeInstances()
	if err != nil {
		return nil, err
	}

	names := make(map[string]bool)

	for _, inst := range instances {
		instName := instanceName(inst)
		if instName != nil {
			names[*instName] = true
		}
	}

	return names, nil
}

func autofilters() []*ec2.Filter {
	values := make([]*string, 0)

	for _, val := range []string{"DAVINCI_ENV_FULL", "DAVINCI_ENV"} {
		values = append(values, aws.String(os.Getenv(val)))
	}

	filters := []*ec2.Filter{
		{
			Name:   aws.String("tag:env"),
			Values: values,
		},
	}
	return filters
}
