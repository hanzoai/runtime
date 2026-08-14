// Copyright 2025 Hanzo Industries Inc.
// SPDX-License-Identifier: AGPL-3.0

package organization

import (
	"github.com/charmbracelet/huh"
	apiclient "github.com/hanzoai/runtime/libs/api-client-go"
	"github.com/hanzoai/runtime/apps/cli/views/common"
)

func GetOrganizationIdFromPrompt(organizationList []apiclient.Organization) (*apiclient.Organization, error) {
	var chosenOrganizationId string
	var organizationOptions []huh.Option[string]

	for _, organization := range organizationList {
		organizationOptions = append(organizationOptions, huh.NewOption(organization.Name, organization.Id))
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose an Organization").
				Options(
					organizationOptions...,
				).
				Value(&chosenOrganizationId),
		).WithTheme(common.GetCustomTheme()),
	)

	if err := form.Run(); err != nil {
		return nil, err
	}

	for _, organization := range organizationList {
		if organization.Id == chosenOrganizationId {
			return &organization, nil
		}
	}

	return nil, nil
}
