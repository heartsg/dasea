
The Dasea authentication and authorization is currently based on
Openstack Keystone service. This doc describes how to install keystone 
service and how to configure it for Dasea project.

1. Installation from Ubuntu Distribution

The general production keystone installation for Ubuntu/Mint is in the following
address:
http://docs.openstack.org/liberty/install-guide-ubuntu/keystone-install.html

For testing purpose, we use KEYSTONE_DBPASS: keystone.
ADMIN_PASS: ADMIN

Note that, keystone is installed with username keystone and group keystone. To
run "keystone-manage db_sync", we have to use 
sudo su -s /bin/sh -c "keystone-manage db_sync" keystone

Also note that, we will need to install pymysql before the above instructions can
be followed:
sudo pip install pymysql

Lastly, the latest keystone verison is liberty, however, the install repo up to 
19/11/2015 is still old kilo version. We have to update apt repository first.

sudo add-apt-repository cloud-archive:liberty

However, this seems doesn't work with linux mint. The cloud-archive only works for
Ubuntu Trusty?

2. Install openstack command line tool

We use python-openstackclient to startup with keystone. The following shows how
to install python-openstackclient from source.

> git clone git://git.openstack.org/openstack/python-openstackclient
> cd python-openstackclient
> virtualenv .env
> source .env/bin/activate
> pip install -r requirements.txt
> python setup.py install

Note that this is latest version of openstack, which may not work together with keystone kilo.
We have to install latest keystone from source or install keystone liberty to make it work.

We can then start using openstack command line tool to create domains, users, projects,
roles and etc. We currently create one super admin user for dasea admin operations. This
user should have all rights to do anything on keystone, mainly used to create/delete domains
and domain/project specific admin users. Once those admin users are created, they shall be
able to create project and project specific users, roles etc. within the domain.

Note that, after database initialized, a domain with name default is created automatically.
We use default domain as super admin user domain.

3. Install keystone from source

Because openstack client and keystone works together only for same version (e.g., the latest
source), it is good to install both keystone and openstack from latest source. We already showed
how to install openstackclient from source, now is keystone source.

We will configure keystone from source, e will then work on how to setup keystone from 
source and let it work together with apache2.

> virtualenv .env
> source .env/bin/activate
> pip install -r requirements.txt
> python setup.py install

Then copy all configuration files into ~/.keystone/, and modify keystone.conf accordingly.





4. Configuring policies for keystone

#create a super domain for super administration purpose
 openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 \
 domain create --description "For super administrators" super
#create a project named super, within super domain
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 project \
 create --domain super --description "Super project" super
#create a user named super, within super domain
#for testing purpose, we make super password super too
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 user \
 create --domain super --password-prompt super
#create a user named admin, within super domain, password is admin
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 user \
 create --domain super --password-prompt admin
#create a role super an assign to the project super and  user super
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 role \
 create super
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 role \
 add --project super --user super super
#create role admin
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 role \
 create admin
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 role \
 add --project super --user admin admin
 

#modify default keystone policy
The default keystone policy (or the overall openstack policy), says that all admins are allowed to
do all operations. However, we want the domain-specific admins only have rights to do parts
of the operations such as create projects/users only for the specific domains. We do not allow admin users
to create domains.

Therefore, we have super role, and only allow super role to do some specific tasks. The modified
keystone policy is given in repo.

#test new created super user and new policy

1. remove admin_token_auth from keystone_paste.ini
2. try to use the new super user and admin user to login
3. super user should be able to create roles, but admin user should not

> openstack --os-auth-url http://localhost:35357/v3 \
  --os-project-domain-name super --os-user-domain-name super \
  --os-project-name super --os-username super --os-auth-type password \
  token issue
  
Record the id field, replace the id with ADMIN in the --os-token above. New roles
can be created/deleted without problem.

> openstack --os-auth-url http://localhost:35357/v3 \
  --os-project-domain-name super --os-user-domain-name super \
  --os-project-name super --os-username admin --os-auth-type password \
  token issue
  
The admin id does not have right to create/delete roles.




