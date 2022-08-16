# GPG

## Create GPG Keys
```
gpg --gen-key --homedir .gnupg
gpg --no-default-keyring --homedir ./.gnupg/ --export-secret-keys > ./.gnupg/secring.gpg
gpg --no-default-keyring --homedir ./.gnupg/ --export > ./.gnupg/pubring.gpg
```

