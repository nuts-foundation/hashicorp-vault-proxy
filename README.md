# hashicorp-vault-proxy
A small proxy server which implements the Nuts Storage API, forwarding calls to a Hashicorp Vault Server.

Running
*******

Build the proxy, then start it and a backing dev Vault server using: 

    $ make docker run

The proxy will be available on port `8210`.

Note: to build the proxy before running the test suite, run:

    $ make docker run

Test suite
**********

To run the test suite that tests compliancy of the proxy with the Nuts Storage API, run:

    $ make run-test

It starts the proxy, Vault and Postman in Docker and runs the test suite.
If the process exits with a non-zero exit code, the test suite failed.
See the Postman output for more information on the failure.

Note: to build the proxy before running the test suite, run:

    $ make docker run-test

Code Generation
***************

Generating code:

To regenerate all code run the ``run-generators`` target from the makefile or use one of the following for a specific group

================ =======================
Group            Command
================ =======================
OpenApi          ``make gen-api``
================ =======================

