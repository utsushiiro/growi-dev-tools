## ベース
[Keycloakを使って、GrowiのSAML認証を使ってみた](https://qiita.com/sahya/items/d8ad66aadcf587c6f0a3)を参考.

## Keycloak側の設定
Keycloak側の設定はkeycloakコンテナに入って`/bin/setup-saml.sh`を実行.

## GROWI側の設定
上記で作成したGROWIレルムの情報は[ここ](http://localhost:8080/auth/realms/GROWI/protocol/saml/descriptor)でXML形式で見ることができる.

GROWI側の設定はdev/env.dev.jsに以下を追加する.
```
  SAML_USES_ONLY_ENV_VARS_FOR_SOME_OPTIONS: true,
  SAML_ENABLED: true,
  SAML_ENTRY_POINT: 'http://localhost:8080/auth/realms/GROWI/protocol/saml',
  SAML_ISSUER: 'localhost:3000', 
  SAML_ATTR_MAPPING_ID: "id",
  SAML_ATTR_MAPPING_USERNAME: "username",
  SAML_ATTR_MAPPING_MAIL: "email",
  SAML_ATTR_MAPPING_FIRST_NAME: "firstName",
  SAML_ATTR_MAPPING_LAST_NAME: "lastName",
  SAML_CERT: '<上記のコマンドの最後に出力された証明書を入れる or 上記のレルムのXMLから抜き出す>'
```

なお, 上記のISSUERはKeycloakのClientのclientIdに対応している.

## よくあるミス

GROWIのログイン画面でSAMLでログインをクリックした後にKeycloakにリダイレクトされたところで「無効なリダイレクト URI です」とでる場合(ログの例は以下の通り).

```
keycloak_1            | 13:46:10,635 WARN  [org.keycloak.events] (default task-15) type=LOGIN_ERROR, realmId=08b1a96e-84e8-4e10-975c-c84cfa7641ef, clientId=null, userId=null, ipAddress=172.18.0.1, error=invalid_redirect_uri
```

この場合にはサイトURLを間違えている場合が多い. アプリ設定 > サイトURL設定で正しいサイトURLを設定する. ポート番号が抜けている等のミスがよくある.
このサイトURLからコールバックURLが作成されてKeycloakへのAuthRequestに利用されるので, ここが間違っていると上記のコマンドで設定したKeycloak側のバリデーション(Valid Redirect URIs)と一致せずにエラーになる.

なお, サイトURLは再起動しないと反映されないため注意が必要.
