## Run benchmark for hashicorp vault

This repo tests whether Vault is able to handle up to 1million secrets
and still return sub milliseconds responses.

Run vault
```
docker run -d --name perf-vault -e VAULT_DEV_ROOT_TOKEN_ID=password -p 8200:8200 hashicorp/vault
```

Run tests
```
go build -o benchvault
./benchvault | tee results.log
```

Stop and remove vault
```bash
docker rm -f perf-vault
```

## Results

**Even at 1m api keys vault is able to serve requests in sub milliseconds.**

Measurements:

```bash
Average response time: 0.560000 ms, errors: 0
Total errors adding 1000000 keys was: 65

Memory: 3.2 GiB
```