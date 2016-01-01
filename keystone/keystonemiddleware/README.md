Keystonemiddleware is based on openstack's keystonemiddleware written in python.
The following docs are copied directly from python's implementation.

Token-based Authentication Middleware

This WSGI component:

* Verifies that incoming client requests have valid tokens by validating
  tokens with the auth service.
* Rejects unauthenticated requests unless the auth_token middleware is in
  ``delay_auth_decision`` mode, which means the final decision is delegated to
  the downstream WSGI component (usually the OpenStack service).
* Collects and forwards identity information based on a valid token
  such as user name, domain, project, etc.

Refer to: http://docs.openstack.org/developer/keystonemiddleware/\
middlewarearchitecture.html


Headers
-------

The auth_token middleware uses headers sent in by the client on the request
and sets headers and environment variables for the downstream WSGI component.

Coming in from initial call from client or customer
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

HTTP_X_AUTH_TOKEN
    The client token being passed in.

HTTP_X_SERVICE_TOKEN
    A service token being passed in.

Used for communication between components
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

WWW-Authenticate
    HTTP header returned to a user indicating which endpoint to use
    to retrieve a new token.

What auth_token adds to the request for use by the OpenStack service
^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

When using composite authentication (a user and service token are
present) additional service headers relating to the service user
will be added. They take the same form as the standard headers but add
``_SERVICE_``. These headers will not exist in the environment if no
service token is present.

HTTP_X_IDENTITY_STATUS, HTTP_X_SERVICE_IDENTITY_STATUS
    Will be set to either ``Confirmed`` or ``Invalid``.

    The underlying service will only see a value of 'Invalid' if the middleware
    is configured to run in ``delay_auth_decision`` mode. As with all such
    headers, ``HTTP_X_SERVICE_IDENTITY_STATUS`` will only exist in the
    environment if a service token is presented. This is different than
    ``HTTP_X_IDENTITY_STATUS`` which is always set even if no user token is
    presented. This allows the underlying service to determine if a
    denial should use ``401 Unauthenticated`` or ``403 Forbidden``.

HTTP_X_DOMAIN_ID, HTTP_X_SERVICE_DOMAIN_ID
    Identity service managed unique identifier, string. Only present if
    this is a domain-scoped token.

HTTP_X_DOMAIN_NAME, HTTP_X_SERVICE_DOMAIN_NAME
    Unique domain name, string. Only present if this is a domain-scoped
    token.

HTTP_X_PROJECT_ID, HTTP_X_SERVICE_PROJECT_ID
    Identity service managed unique identifier, string. Only present if
    this is a project-scoped token.

HTTP_X_PROJECT_NAME, HTTP_X_SERVICE_PROJECT_NAME
    Project name, unique within owning domain, string. Only present if
    this is a project-scoped token.

HTTP_X_PROJECT_DOMAIN_ID, HTTP_X_SERVICE_PROJECT_DOMAIN_ID
    Identity service managed unique identifier of owning domain of
    project, string.  Only present if this is a project-scoped v3 token. If
    this variable is set, this indicates that the PROJECT_NAME can only
    be assumed to be unique within this domain.

HTTP_X_PROJECT_DOMAIN_NAME, HTTP_X_SERVICE_PROJECT_DOMAIN_NAME
    Name of owning domain of project, string. Only present if this is a
    project-scoped v3 token. If this variable is set, this indicates that
    the PROJECT_NAME can only be assumed to be unique within this domain.

HTTP_X_USER_ID, HTTP_X_SERVICE_USER_ID
    Identity-service managed unique identifier, string.

HTTP_X_USER_NAME, HTTP_X_SERVICE_USER_NAME
    User identifier, unique within owning domain, string.

HTTP_X_USER_DOMAIN_ID, HTTP_X_SERVICE_USER_DOMAIN_ID
    Identity service managed unique identifier of owning domain of
    user, string. If this variable is set, this indicates that the USER_NAME
    can only be assumed to be unique within this domain.

HTTP_X_USER_DOMAIN_NAME, HTTP_X_SERVICE_USER_DOMAIN_NAME
    Name of owning domain of user, string. If this variable is set, this
    indicates that the USER_NAME can only be assumed to be unique within
    this domain.

HTTP_X_ROLES, HTTP_X_SERVICE_ROLES
    Comma delimited list of case-sensitive role names.

HTTP_X_SERVICE_CATALOG
    service catalog (optional, JSON string).

    For compatibility reasons this catalog will always be in the V2 catalog
    format even if it is a v3 token.

    .. note:: This is an exception in that it contains 'SERVICE' but relates to
        a user token, not a service token. The existing user's catalog can be
        very large; it was decided not to present a catalog relating to the
        service token to avoid using more HTTP header space.

HTTP_X_TENANT_ID
    *Deprecated* in favor of HTTP_X_PROJECT_ID.

    Identity service managed unique identifier, string. For v3 tokens, this
    will be set to the same value as HTTP_X_PROJECT_ID.

HTTP_X_TENANT_NAME
    *Deprecated* in favor of HTTP_X_PROJECT_NAME.

    Project identifier, unique within owning domain, string. For v3 tokens,
    this will be set to the same value as HTTP_X_PROJECT_NAME.

HTTP_X_TENANT
    *Deprecated* in favor of HTTP_X_TENANT_ID and HTTP_X_TENANT_NAME.

    Identity server-assigned unique identifier, string. For v3 tokens, this
    will be set to the same value as HTTP_X_PROJECT_ID.

HTTP_X_USER
    *Deprecated* in favor of HTTP_X_USER_ID and HTTP_X_USER_NAME.

    User name, unique within owning domain, string.

HTTP_X_ROLE
    *Deprecated* in favor of HTTP_X_ROLES.

    Will contain the same values as HTTP_X_ROLES.

Environment Variables
^^^^^^^^^^^^^^^^^^^^^

These variables are set in the request environment for use by the downstream
WSGI component.

keystone.token_info
    Information about the token discovered in the process of validation.  This
    may include extended information returned by the token validation call, as
    well as basic information about the project and user.

keystone.token_auth
    A keystoneclient auth plugin that may be used with a
    :py:class:`keystoneclient.session.Session`. This plugin will load the
    authentication data provided to auth_token middleware.


Configuration
-------------

auth_token middleware configuration can be in the main application's
configuration file, e.g. in ``nova.conf``:

.. code-block:: ini

  [keystone_authtoken]
  auth_plugin = password
  auth_url = http://keystone:35357/
  username = nova
  user_domain_id = default
  password = whyarewestillusingpasswords
  project_name = service
  project_domain_id = default

Configuration can also be in the ``api-paste.ini`` file with the same options,
but this is discouraged.

Swift
-----

When deploy auth_token middleware with Swift, user may elect to use Swift
memcache instead of the local auth_token memcache. Swift memcache is passed in
from the request environment and it's identified by the ``swift.cache`` key.
However it could be different, depending on deployment. To use Swift memcache,
you must set the ``cache`` option to the environment key where the Swift cache
object is stored.
