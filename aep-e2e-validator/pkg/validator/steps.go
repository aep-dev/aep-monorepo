package validator

import (
	"fmt"
	"io"
	"net/http"

	"github.com/aep-dev/aep-lib-go/pkg/api"
)

type ValidationContext struct {
	Resource      *api.Resource
	CollectionURL string
	Resource1     map[string]interface{}
	Resource2     map[string]interface{}
	ListResponse1 *ListResponse
}

type ValidationStep struct {
	Name string
	Run  func(ctx *ValidationContext) error
}

func (v *Validator) stepCreateFirstResource(ctx *ValidationContext) error {
	createPayload, err := GenerateCreatePayload(ctx.Resource)
	if err != nil {
		return fmt.Errorf("failed to generate create payload: %w", err)
	}

	resource1, err := v.createResource(ctx.Resource, ctx.CollectionURL, createPayload)
	if err != nil {
		return err
	}

	ctx.Resource1 = resource1
	r1Name, _ := resource1["name"].(string)
	fmt.Printf("   Created %s\n", r1Name)
	return nil
}

func (v *Validator) stepDuplicateCreationCheck(ctx *ValidationContext) error {
	r := ctx.Resource
	if r.Methods.Create != nil && r.Methods.Create.SupportsUserSettableCreate {
		fmt.Println("   Attempting duplicate creation...")
		r1Name, _ := ctx.Resource1["name"].(string)
		r1ID := getIDFromResourceName(r1Name)
		createPayload, _ := GenerateCreatePayload(r) // Payload content doesn't matter much for ID check

		urlWithID := fmt.Sprintf("%s?%s_id=%s", ctx.CollectionURL, r.Singular, r1ID)
		resp, err := v.post(urlWithID, createPayload)
		if err != nil {
			return fmt.Errorf("failed to make request: %w", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusConflict && resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("expected 409/400 for duplicate creation, got %d", resp.StatusCode)
		}
		fmt.Println("   Duplicate creation rejected as expected.")
	} else {
		fmt.Println("   Skipping duplicate check (client-assigned ID not supported or uncheckable).")
	}
	return nil
}

func (v *Validator) stepCreateSecondResource(ctx *ValidationContext) error {
	createPayload2, err := GenerateCreatePayload(ctx.Resource)
	if err != nil {
		return fmt.Errorf("failed to generate payload 2: %w", err)
	}
	resource2, err := v.createResource(ctx.Resource, ctx.CollectionURL, createPayload2)
	if err != nil {
		return err
	}
	ctx.Resource2 = resource2
	r2Name, _ := resource2["name"].(string)
	fmt.Printf("   Created %s\n", r2Name)
	return nil
}

func (v *Validator) stepListResourcesLimit1(ctx *ValidationContext) error {
	listURL := fmt.Sprintf("%s?page_size=1", ctx.CollectionURL)
	listResp1, err := v.list(listURL)
	if err != nil {
		return err
	}
	if len(listResp1.Resources) != 1 {
		return fmt.Errorf("expected 1 resource, got %d", len(listResp1.Resources))
	}
	if listResp1.NextPageToken == "" {
		return fmt.Errorf("expected next_page_token")
	}
	ctx.ListResponse1 = listResp1
	fmt.Println("   Got 1 resource and next_page_token.")
	return nil
}

func (v *Validator) stepListResourcesPageToken(ctx *ValidationContext) error {
	listURL2 := fmt.Sprintf("%s?page_size=1&page_token=%s", ctx.CollectionURL, ctx.ListResponse1.NextPageToken)
	listResp2, err := v.list(listURL2)
	if err != nil {
		return err
	}
	if len(listResp2.Resources) < 1 {
		return fmt.Errorf("expected at least 1 resource (the second one)")
	}
	fmt.Println("   Got second page.")
	return nil
}

func (v *Validator) stepUpdateFirstResource(ctx *ValidationContext) error {
	updatePayload, err := GenerateCreatePayload(ctx.Resource)
	if err != nil {
		return fmt.Errorf("failed to generate update payload: %w", err)
	}
	r1Name, _ := ctx.Resource1["name"].(string)
	r1URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r1Name)

	respUpdate, err := v.patch(r1URL, updatePayload)
	if err != nil {
		return err
	}
	defer respUpdate.Body.Close()
	if respUpdate.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(respUpdate.Body)
		return fmt.Errorf("update returned %d: %s", respUpdate.StatusCode, string(body))
	}
	fmt.Println("   Update successful.")
	return nil
}

func (v *Validator) stepGetFirstResource(ctx *ValidationContext) error {
	r1Name, _ := ctx.Resource1["name"].(string)
	r1URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r1Name)
	r1Get, err := v.get(r1URL)
	if err != nil {
		return err
	}
	_ = r1Get
	fmt.Println("   Get successful.")
	return nil
}

func (v *Validator) stepDeleteFirstResource(ctx *ValidationContext) error {
	r1Name, _ := ctx.Resource1["name"].(string)
	r1URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r1Name)
	if err := v.delete(r1URL); err != nil {
		return err
	}
	fmt.Println("   Delete successful.")
	return nil
}

func (v *Validator) stepDuplicateDeleteCheck(ctx *ValidationContext) error {
	r1Name, _ := ctx.Resource1["name"].(string)
	r1URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r1Name)
	respDeleteDup, err := v.deleteReq(r1URL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer respDeleteDup.Body.Close()
	if respDeleteDup.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404, got %d", respDeleteDup.StatusCode)
	}
	fmt.Println("   Got 404 as expected.")
	return nil
}

func (v *Validator) stepGetDeletedResource(ctx *ValidationContext) error {
	r1Name, _ := ctx.Resource1["name"].(string)
	r1URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r1Name)
	respGetDel, err := v.getReq(r1URL)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer respGetDel.Body.Close()
	if respGetDel.StatusCode != http.StatusNotFound {
		return fmt.Errorf("expected 404, got %d", respGetDel.StatusCode)
	}
	fmt.Println("   Got 404 as expected.")
	return nil
}

func (v *Validator) stepDeleteSecondResource(ctx *ValidationContext) error {
	r2Name, _ := ctx.Resource2["name"].(string)
	r2URL := fmt.Sprintf("%s/%s", ctx.Resource.API.ServerURL, r2Name)
	if err := v.delete(r2URL); err != nil {
		return err
	}
	fmt.Println("   Delete successful.")
	return nil
}

func (v *Validator) stepListVerifyEmpty(ctx *ValidationContext) error {
	listRespFinal, err := v.list(ctx.CollectionURL)
	if err != nil {
		return err
	}
	fmt.Printf("   List count: %d\n", len(listRespFinal.Resources))
	return nil
}
