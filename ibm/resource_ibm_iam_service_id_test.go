// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/IBM-Cloud/bluemix-go/models"

	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIBMIAMServiceID_Basic(t *testing.T) {
	var conf models.ServiceID
	name := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	updateName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIAMServiceIDDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIAMServiceIDBasic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMServiceIDExists("ibm_iam_service_id.serviceID", conf),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "name", name),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "tags.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIAMServiceIDUpdateWithSameName(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMServiceIDExists("ibm_iam_service_id.serviceID", conf),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "name", name),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "description", "ServiceID for test scenario1"),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "tags.#", "3"),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIAMServiceIDUpdate(updateName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "name", updateName),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "description", "ServiceID for test scenario2"),
					resource.TestCheckResourceAttr("ibm_iam_service_id.serviceID", "tags.#", "1"),
				),
			},
		},
	})
}

func TestAccIBMIAMServiceID_import(t *testing.T) {
	var conf models.ServiceID
	name := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))
	resourceName := "ibm_iam_service_id.serviceID"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIAMServiceIDDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIAMServiceIDTag(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIAMServiceIDExists(resourceName, conf),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "ServiceID for test scenario2"),
				),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIAMServiceIDDestroy(s *terraform.State) error {
	rsContClient, err := testAccProvider.Meta().(ClientSession).IAMAPI()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_iam_service_id" {
			continue
		}

		serviceIDUUID := rs.Primary.ID

		// Try to find the key
		_, err := rsContClient.ServiceIds().Get(serviceIDUUID)

		if err == nil {
			return fmt.Errorf("ServiceID still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error waiting for serviceID (%s) to be destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckIBMIAMServiceIDExists(n string, obj models.ServiceID) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		rsContClient, err := testAccProvider.Meta().(ClientSession).IAMAPI()
		if err != nil {
			return err
		}
		serviceIDUUID := rs.Primary.ID

		serviceID, err := rsContClient.ServiceIds().Get(serviceIDUUID)

		if err != nil {
			return err
		}

		obj = serviceID
		return nil
	}
}

func testAccCheckIBMIAMServiceIDBasic(name string) string {
	return fmt.Sprintf(`
		
		resource "ibm_iam_service_id" "serviceID" {
			name = "%s"
			tags = ["tag1", "tag2"]
	  	}
	`, name)
}

func testAccCheckIBMIAMServiceIDUpdateWithSameName(name string) string {
	return fmt.Sprintf(`
		
		resource "ibm_iam_service_id" "serviceID" {
			name        = "%s"
			description = "ServiceID for test scenario1"
			tags        = ["tag1", "tag2", "db"]
	  	}
	`, name)
}

func testAccCheckIBMIAMServiceIDUpdate(updateName string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_service_id" "serviceID" {
			name              = "%s"		
			description       = "ServiceID for test scenario2"
			tags              = ["tag1"]
		}
	`, updateName)
}

func testAccCheckIBMIAMServiceIDTag(name string) string {
	return fmt.Sprintf(`

		resource "ibm_iam_service_id" "serviceID" {
			name              = "%s"		
			description       = "ServiceID for test scenario2"
		}
	`, name)
}
