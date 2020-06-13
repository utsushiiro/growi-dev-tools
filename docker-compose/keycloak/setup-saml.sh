#!/bin/bash -e

PATH=$PATH:$JBOSS_HOME/bin
KEYCLOAK_URL="http://keycloak:8080"
REALM_NAME="GROWI"
KEYCLOAK_USER="admin"
KEYCLOAK_PASSWORD="password"
GROWI_APP_FQDN="localhost:3000"
HACKMD_FQDN="localhost:3010"

# Keycloakにログイン
kcadm.sh config credentials --server ${KEYCLOAK_URL}/auth --realm master --user ${KEYCLOAK_USER} --password ${KEYCLOAK_PASSWORD}

# GROWI用のレルムを作成する(存在しない場合のみ)
kcadm.sh get realms/${REALM_NAME} > /dev/null || kcadm.sh create realms -f - << EOF
{ "realm": "${REALM_NAME}", "enabled": true, "internationalizationEnabled": true, "defaultLocale": "ja" }
EOF

# GROWI用のクライアントプロファイルを作成する(存在しない場合のみ)
kcadm.sh get clients -r ${REALM_NAME} --fields=clientId | grep ${GROWI_APP_FQDN} > /dev/null 2>&1 || kcadm.sh create clients -r ${REALM_NAME} -f - << EOF
{
  "clientId": "${GROWI_APP_FQDN}",
  "name": "GROWI",
  "rootUrl": "http://${GROWI_APP_FQDN}/",
  "enabled": true,
  "redirectUris": [
    "http://${GROWI_APP_FQDN}/*"
  ],
  "protocol": "saml",
  "attributes": {
    "saml.authnstatement": "true",
    "saml.assertion.signature": "true",
    "saml.force.post.binding": "true",
    "saml.server.signature": "true",
    "saml.signature.algorithm": "RSA_SHA1",
    "saml.client.signature" : "false",
    "saml_force_name_id_format": "true",
    "saml_name_id_format": "username",
    "saml_assertion_consumer_url_redirect": "http://${GROWI_APP_FQDN}/passport/saml/callback",
    "saml_assertion_consumer_url_post": "http://${GROWI_APP_FQDN}/",
    "saml_single_logout_service_url_redirect": "http://${GROWI_APP_FQDN}/passport/saml/callback"
  },
  "protocolMappers": [
    {
      "name": "username",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "username",
        "friendly.name": "username",
        "attribute.name": "username"
      }
    },
    {
      "name": "id",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "id",
        "friendly.name": "id",
        "attribute.name": "id"
      }
    },
    {
      "name": "lastname",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "lastname",
        "friendly.name": "lastname",
        "attribute.name": "lastname"
      }
    },
    {
      "name": "email",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "email",
        "friendly.name": "email",
        "attribute.name": "email"
      }
    },
    {
      "name": "firstname",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "firstname",
        "friendly.name": "firstname",
        "attribute.name": "firstname"
      }
    }
  ]
}
EOF

# HackMD用のクライアントプロファイルを作成する(存在しない場合のみ)
kcadm.sh get clients -r ${REALM_NAME} --fields=clientId | grep ${HACKMD_FQDN} > /dev/null 2>&1 || kcadm.sh create clients -r ${REALM_NAME} -f - << EOF
{
  "clientId": "${HACKMD_FQDN}",
  "name": "HackMD",
  "rootUrl": "http://${HACKMD_FQDN}/",
  "enabled": true,
  "redirectUris": [
    "http://${HACKMD_FQDN}/*"
  ],
  "protocol": "saml",
  "attributes": {
    "saml.authnstatement": "true",
    "saml.assertion.signature": "true",
    "saml.force.post.binding": "true",
    "saml.server.signature": "true",
    "saml.signature.algorithm": "RSA_SHA1",
    "saml.client.signature" : "false",
    "saml_force_name_id_format": "true",
    "saml_name_id_format": "username",
    "saml_assertion_consumer_url_redirect": "http://${HACKMD_FQDN}/auth/saml/callback",
    "saml_assertion_consumer_url_post": "http://${HACKMD_FQDN}/",
    "saml_single_logout_service_url_redirect": "http://${HACKMD_FQDN}/auth/saml/callback"
  },
  "protocolMappers": [
    {
      "name": "username",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "username",
        "friendly.name": "username",
        "attribute.name": "username"
      }
    },
    {
      "name": "id",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "id",
        "friendly.name": "id",
        "attribute.name": "id"
      }
    },
    {
      "name": "email",
      "protocol": "saml",
      "protocolMapper": "saml-user-property-mapper",
      "consentRequired": false,
      "config": {
        "user.attribute": "email",
        "friendly.name": "email",
        "attribute.name": "email"
      }
    }
  ]
}
EOF

CERTIFICATE=`kcadm.sh get keys -r ${REALM_NAME} | jq '.keys[] | select(.type == "RSA") | .certificate' -r`
echo $CERTIFICATE