// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START iam_quickstartv2]

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/cloudresourcemanager/v1"
)

func main() {

	// TODO: Add your project ID
	projectID := flag.String("project_id", "", "Cloud Project ID")
	// TODO: Add the ID of your member in the form "user:member@example.com"
	member := flag.String("member_id", "", "Your member ID")
	flag.Parse()

	// The role to be granted
	var role string = "roles/logging.logWriter"

	// Initializes the Cloud Resource Manager service
	crmService := initializeService()

	// Grants your member the "Log writer" role for your project
	addBinding(crmService, *projectID, *member, role)

	// Gets the project's policy and prints all members with the "Log Writer" role
	policy := getPolicy(crmService, *projectID)
	var binding *cloudresourcemanager.Binding = nil
	bindings := policy.Bindings
	for b := range bindings {
		if bindings[b].Role == role {
			binding = bindings[b]
			break
		}
	}
	fmt.Println("Role: ", binding.Role)
	fmt.Print("Members: ")
	for m := range binding.Members {
		fmt.Print("[", binding.Members[m], "] ")
	}

	// Removes member from the "Log writer" role
	removeMember(crmService, *projectID, *member, role)

}

// initializeService initializes a new Cloud Resource Manager service
func initializeService() *cloudresourcemanager.Service {

	ctx := context.Background()
	crmService, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		log.Fatalf("cloudresourcemanager.NewService: %v", err)
	}
	return crmService

}

// addBinding adds the member to the project's IAM policy
func addBinding(crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Finds the role binding in the policy, if it exists
	bindings := policy.Bindings
	var binding *cloudresourcemanager.Binding = nil
	for b := range bindings {
		if bindings[b].Role == role {
			binding = bindings[b]
			break
		}
	}

	if binding != nil {
		// If the binding exists, adds the member to the binding
		binding.Members = append(binding.Members, member)
	} else {
		// If the binding does not exist, adds a new binding to the policy
		binding = new(cloudresourcemanager.Binding)
		binding.Role = role
		binding.Members = []string{member}
		policy.Bindings = append(policy.Bindings, binding)
	}

	setPolicy(crmService, projectID, policy)

}

// removeMember removes the member from the project's IAM policy
func removeMember(crmService *cloudresourcemanager.Service, projectID, member, role string) {

	policy := getPolicy(crmService, projectID)

	// Finds the binding in the policy
	bindings := policy.Bindings
	var binding *cloudresourcemanager.Binding = nil
	var bindingIndex int
	for b := range bindings {
		if bindings[b].Role == role {
			binding = bindings[b]
			bindingIndex = b
			break
		}
	}

	// Order doesn't matter for bindings or members, so to remove, move the last item
	// into the removed spot and shrink the slice.
	if len(binding.Members) == 1 {
		// If the member is the only member in the binding, removes the binding
		last := len(bindings) - 1
		bindings[bindingIndex] = bindings[last]
		bindings[last] = nil
		policy.Bindings = bindings[:last]
	} else {
		// If there is more than one member in the binding, removes the member
		var memberIndex int
		for i, mm := range binding.Members {
			if mm == member {
				memberIndex = i
			}
		}
		last := len(bindings[bindingIndex].Members) - 1
		binding.Members[memberIndex] = binding.Members[last]
		binding.Members[last] = ""
		binding.Members = binding.Members[:last]
	}

	setPolicy(crmService, projectID, policy)

}

// getPolicy gets the project's IAM policy
func getPolicy(crmService *cloudresourcemanager.Service, projectID string) *cloudresourcemanager.Policy {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.GetIamPolicyRequest)
	policy, err := crmService.Projects.GetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.GetIamPolicy: %v", err)
	}

	return policy
}

// setPolicy sets the project's IAM policy
func setPolicy(crmService *cloudresourcemanager.Service, projectID string, policy *cloudresourcemanager.Policy) {

	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	request := new(cloudresourcemanager.SetIamPolicyRequest)
	request.Policy = policy
	policy, err := crmService.Projects.SetIamPolicy(projectID, request).Do()
	if err != nil {
		log.Fatalf("Projects.SetIamPolicy: %v", err)
	}
}

// [END iam_quickstartv2]
