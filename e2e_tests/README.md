# End-to-end tests

This set of tests is intended for functionality validation on an existing WAPI
instance.

## How to run E2E tests
 
1. Export parameters of accessible and running `WAPI` instance:
   ```bash
   export INFOBLOX_SERVER=<WAPI HOST IP> INFOBLOX_USERNAME=<WAPI USERNAME> INFOBLOX_PASSWORD=<WAPI PASSWORD> WAPI_VERSION=<WAPI VERSION>
   ```

2. You can use `go test` utility to run the tests in the `e2e_tests` directory:
   ```bash
   go test -v ./e2e_tests
   ```

3. Instead of `go test` you can also use the [`ginkgo`](https://onsi.github.io/ginkgo/#installing-ginkgo) utility, 
   which supports test filtering by label, and lots of [additional features](https://onsi.github.io/ginkgo/#running-specs):
   ```bash
   # Run all e2e tests
   ginkgo e2e_tests
   # Run only read-only e2e tests
   ginkgo --label-filter=RO e2e_tests
   ```

## Warning

Please don't run those tests on the production WAPI instance.
