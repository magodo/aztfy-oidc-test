terraform {
  required_providers {
    restful = {
      source  = "magodo/restful"
      version = "0.16.1"
    }
  }
}

variable "id_token_request_url" {
  type = string
}

variable "id_token_request_token" {
  type = string
}

variable "tenant_id" {
  type = string
}

variable "client_id" {
  type = string
}

variable "sub_id" {
  type = string
}

provider "restful" {
  base_url = var.id_token_request_url

  security = {
    http = {
      token = {
        token = var.id_token_request_token
      }
    }
  }

  alias = "id_token"
}

data "restful_resource" "id_token" {
  id = ""
  query = {
    "audience" : ["api://AzureADTokenExchange"]
  }

  provider = restful.id_token
}

provider "restful" {
  base_url = "https://login.microsoftonline.com"
  alias    = "access_token"
}

resource "restful_operation" "access_token" {
  path   = "/${var.tenant_id}/oauth2/v2.0/token"
  method = "POST"
  header = {
    Accept       = "application/json"
    Content-Type = "application/x-www-form-urlencoded"
  }
  body = {
    client_assertion      = data.restful_resource.id_token.output.value
    client_assertion_type = "urn:ietf:params:oauth:client-assertion-type:jwt-bearer"
    client_id             = var.client_id
    grant_type            = "client_credentials"
    scope                 = "https://management.azure.com/.default"
  }
  provider = restful.access_token
}

provider "restful" {
  base_url = "https://management.azure.com"
  security = {
    http = {
      token = {
        token = restful_operation.access_token.output.access_token
      }
    }
  }
}

data "restful_resource" "test" {
  id = "/subscriptions/var.sub_id/resourceGroups/aztfy"
  query = {
    api-version = ["2020-06-01"]
  }
  depends_on = [restful_operation.access_token]
}
