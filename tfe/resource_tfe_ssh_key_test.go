package tfe

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	tfe "github.com/hashicorp/go-tfe"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTFESSHKey_basic(t *testing.T) {
	sshKey := &tfe.SSHKey{}
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFESSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTFESSHKey_basic, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESSHKeyExists(
						"tfe_ssh_key.foobar", sshKey),
					testAccCheckTFESSHKeyAttributes(sshKey),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "name", "ssh-key-test"),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "key", "SSH-KEY-CONTENT"),
				),
			},
		},
	})
}

func TestAccTFESSHKey_update(t *testing.T) {
	sshKey := &tfe.SSHKey{}
	rInt := rand.New(rand.NewSource(time.Now().UnixNano())).Int()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckTFESSHKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccTFESSHKey_basic, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESSHKeyExists(
						"tfe_ssh_key.foobar", sshKey),
					testAccCheckTFESSHKeyAttributes(sshKey),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "name", "ssh-key-test"),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "key", "SSH-KEY-CONTENT"),
				),
			},

			{
				Config: fmt.Sprintf(testAccTFESSHKey_update, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTFESSHKeyExists(
						"tfe_ssh_key.foobar", sshKey),
					testAccCheckTFESSHKeyAttributesUpdated(sshKey),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "name", "ssh-key-updated"),
					resource.TestCheckResourceAttr(
						"tfe_ssh_key.foobar", "key", "UPDATED-SSH-KEY-CONTENT"),
				),
			},
		},
	})
}

func testAccCheckTFESSHKeyExists(
	n string, sshKey *tfe.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tfeClient := testAccProvider.Meta().(*tfe.Client)

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		sk, err := tfeClient.SSHKeys.Read(ctx, rs.Primary.ID)
		if err != nil {
			return err
		}

		if sk == nil {
			return fmt.Errorf("SSH key not found")
		}

		*sshKey = *sk

		return nil
	}
}

func testAccCheckTFESSHKeyAttributes(
	sshKey *tfe.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sshKey.Name != "ssh-key-test" {
			return fmt.Errorf("Bad name: %s", sshKey.Name)
		}
		return nil
	}
}

func testAccCheckTFESSHKeyAttributesUpdated(
	sshKey *tfe.SSHKey) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sshKey.Name != "ssh-key-updated" {
			return fmt.Errorf("Bad name: %s", sshKey.Name)
		}
		return nil
	}
}

func testAccCheckTFESSHKeyDestroy(s *terraform.State) error {
	tfeClient := testAccProvider.Meta().(*tfe.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "tfe_ssh_key" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		_, err := tfeClient.SSHKeys.Read(ctx, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("SSH key %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

const testAccTFESSHKey_basic = `
resource "tfe_organization" "foobar" {
  name  = "tst-terraform-%d"
  email = "admin@company.com"
}

resource "tfe_ssh_key" "foobar" {
  name         = "ssh-key-test"
  organization = "${tfe_organization.foobar.id}"
  key          = "SSH-KEY-CONTENT"
}`

const testAccTFESSHKey_update = `
resource "tfe_organization" "foobar" {
  name  = "tst-terraform-%d"
  email = "admin@company.com"
}

resource "tfe_ssh_key" "foobar" {
  name         = "ssh-key-updated"
  organization = "${tfe_organization.foobar.id}"
  key          = "UPDATED-SSH-KEY-CONTENT"
}`
