
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

We will configure keystone from source and let it work together with apache2.

Too trouble some, try without apache first.

4. Configuring policies for keystone

#create a project named admin, within default domain
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 project \
 create --domain default --description "Admin project" admin
#create a user named super, within project admin and default domain
openstack --os-token ADMIN --os-url http://localhost:35357/v3 --os-identity-api-version 3 user \
 create --domain default super
#create a role super an assign to the user super
 


