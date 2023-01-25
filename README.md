# hashicorp-vault-proxy
A small proxy server which implements the Nuts Storage API, forwarding calls to a Hashicorp Vault Server.

Running
*******

To build the application and start it with a Vault server, run:

    $ make build start

The proxy will be available on port `8210`. The Vault server will run in development mode.

To stop the services, run:

    $ make stop

To reset the services, effectively removing the Docker containers and volumes (including the stored private keys), run:

    $ make reset

Test suite
**********

To run the test suite that tests compliance of the proxy with the Nuts Storage API, run:

    $ make api-test

It starts the proxy, Vault and Postman in Docker and runs the test suite.
If the process exits with a non-zero exit code, the test suite failed.
See the Postman output for more information on the failure.

Note: to build the proxy before running the test suite, run:

    $ make build api-test

Code Generation
***************

Generating code:

To regenerate all code run the ``run-generators`` target:

    $ make run-generators

