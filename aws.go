package main

import (
	"fmt"
	"time"

	//"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strings"
)

var imageMap = make(map[string]string)
var ec2Svc = ec2.New(getSession())

// InstanceData is used to figure out
// where to SSH to
type InstanceData struct {
	PrivateIP  string
	InstanceID string
	Name       string
	Az         string
	Os         string
	LaunchTime time.Time
}

func getSession() *session.Session {
	return session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
}

func tagValue(tags []*ec2.Tag, name string) string {
	for _, v := range tags {
		if *v.Key == name {
			return *v.Value
		}
	}
	return ""
}

func getInstanceData() []*InstanceData {

	var instances []*InstanceData
	result, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		fmt.Println("Error in getInstanceData", err)
		panic(nil)
	}
	for _, reservation := range result.Reservations {
		awsData := reservation.Instances[0]
		data := InstanceData{
			PrivateIP:  *awsData.NetworkInterfaces[0].PrivateIpAddress,
			Name:       tagValue(awsData.Tags, "Name"),
			InstanceID: *awsData.InstanceId,
			Az:         *awsData.Placement.AvailabilityZone,
			Os:         getAmiOS(awsData.ImageId),
			LaunchTime: *awsData.LaunchTime,
		}
		instances = append(instances, &data)
	}
	return instances
}

func getAmiOS(amiID *string) string {

	cachedImage, ok := imageMap[*amiID]
	if ok {
		return cachedImage
	}
	imageFilter := "image-id"
	filter := &ec2.Filter{
		Name:   &imageFilter,
		Values: []*string{amiID},
	}
	opts := &ec2.DescribeImagesInput{
		Filters: []*ec2.Filter{filter},
	}
	response, err := ec2Svc.DescribeImages(opts)
	if err != nil {
		fmt.Println("AWS ERROR: ", err)
	}
	imageMap[*amiID] = *response.Images[0].Name
	return *response.Images[0].Name
}

// Looks at the AMI id and figures out the default User
func getDefaultUser(amiName string) string {
	switch {
	case strings.Contains(amiName, "ubuntu"): // Ubuntu
		return "ubuntu"
	default: // The rest?
		return "ec2-user"
	}
}
