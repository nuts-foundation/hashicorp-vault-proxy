# hashicorp-vault-proxy
A small proxy server which implements the Nuts Storage API, forwarding calls to a Hashicorp Vault Server.

Running
*******

Build the proxy, then start it and a backing dev Vault server using: 

    $ make docker run

The proxy will be available on port `8210`.

Code Generation
***************

Generating code:

To regenerate all code run the ``run-generators`` target from the makefile or use one of the following for a specific group

================ =======================
Group            Command
================ =======================
OpenApi          ``make gen-api``
================ =======================

