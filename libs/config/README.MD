# Configuration Library

## Remote Configuration
Remote configuration is abstraction of get configuration on remote(vault,firebase,db, etc) source

### Vault
[Vault](https://www.vaultproject.io/docs/what-is-vault#what-is-vault) is an identity-based secrets and encryption management system. A secret is anything that you want to tightly control access to, such as API encryption keys, passwords, or certificates.
Your can use vault as remote config storage

## Note
Currently only support  vault for remote configuration

## Example

- [Abstract](https://bitbucket.org/efishery/go-efishery/src/master/libs/config/example/abstract/main.go)
- [Your Own Vault](https://bitbucket.org/efishery/go-efishery/src/master/libs/config/example/vault/main.go)

# Milestone
[x] Abstraction remote config for global call

[x] Vault remote config for storage

[x] Vault auto renew token

[x] Cache response data wrapper
