# ocm_aws

The is packaged aws functions to be used for the API or ROSA automation testing
There are two clients defined in pkg/aws_client: aws_v1 and aws_v2
 package aws functions in aws_client package.

 Call function for enhance implementation in ocm/aws package

 How to prepare VPC:
    vpcCidr := "11.0.0.0/16"
	vpc, err := vpc.PrepareVPC("xueli-t", "us-west-2", vpcCidr, true, "a", "b", "c")
	fmt.Println(err)

How to delete vpc and the reourses:
	vpcID := vpc.VpcID
	err = vpc.DeleteVPCChain(vpcID)

# Build a binary for the command call
* Call below command to install the ocmqe package

	`$ go install github.com/openshift-qe/openshift-rosa-cli/ocmqe@latest`

* Make the go binary path to be in PATH

	`$ export PATH=$PATH:/$GOPATH/bin:/usr/local/bin`

# How to use the binary of ocmqe to create resources
* Check the help message

	`$ ocmqe -h`

* Create a vpc on the indicated region
	* When you have a vpc on the indicated region and you want to re-use it, if cannot find create it

	`$ ocmqe create vpc --region us-west-2 --name <your-alias>-vpc`

	`$ ocmqe create vpc --region us-west-2 --name <your-alias>-vpc --find-existing`

* Prepare a pair of subnets on the indicated zone of the region. 
	* NOTE: If the subnets had been existing, it will reuse the exsiting subnets in the zone

	`$ ocmqe create subnets --region us-west-2 --zones a --vpc-id <vpc id>`
    

* Prepare proxy server
    * Parameters Description
      * --region: Specifies the region where your VPC is located.
      * --vpc-id: The id of the VPC under which your MITM proxy will be deployed.
      * --zone: Relates to the region and the subnet you have created. It is used to locate an existing public subnet or create a new one if none exists.
      * --ca-file: The file path where the MITM proxy’s generated root CA will be stored. Clients will use this root CA to establish secure communication through the MITM proxy.
      * --keypair-name: The name of the key pair. Ensure this name does not already exist.
      * --privatekey-path: The directory where the SSH credential file will be stored. This file enables users to SSH into the VM where the MITM proxy is installed.
        
	
	`$ ocmqe create proxy --region us-west-2 --vpc-id <vpc id> --zone <e.g. us-west-2a>  --ca-file <ca file path, eg. ~/proxy/ca-file.pem> --keypair-name <keypair name> --privatekey-path <the directory which store privatekey, e.g. ~/proxy/>`
    

* Prepare addtional security groups, output the sg IDs with comma seperated

	`$ ocmqe create security-groups --region us-west-2 --vpc-id <vpc id> --count <security group number>`

* Tag a resource

	`$ ocmqe tag --resource-id <resource id> --region us-west-2 --tags aaa:ddd,fff:bbb`

* Delete a tag from a resource

	`$ ocmqe delete tag --resource-id <resource id> --region us-west-2 --tag-key <key> --tag-value <value>`

* Clean the vpc and the resources, there is a flag --total-clean supported to do a total clean even the resources is not created by this package

	`$ ocmqe delete vpc --vpc-id <vpc id> --region us-west-2`
