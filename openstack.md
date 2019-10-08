# Instalacion en OpenStack

### 1 - Creacion de instancia
Flavor: CPU1RAM2048HD10GB
Image: RHEL7.6-V1.2

### 2 - Instalacion de artefactos

~~~bash
wget https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-infraestructura-raw/oracle/instantclient/oracle-instantclient19.3-basic-19.3.0.0.0-1.x86_64.rpm
sudo yum install -y oracle-instantclient19.3-basic-19.3.0.0.0-1.x86_64.rpm
wget https://nexus.cloudint.afip.gob.ar/nexus/repository/fisca-infraestructura-raw/mflow/1.0.0/mflow-1.0.0.tar.gz
tar xfvz mflow-1.0.0.tar.gz
sudo mv mflow /opt/
~~~

### 3 - Seteo de variables de entorno

~~~bash
echo 'LD_LIBRARY_PATH=/usr/lib/oracle/19.3/client64/lib:$LD_LIBRARY_PATH' >> ~/.bashrc
source ~/.bashrc
~~~

### 4 - Configuracion de datasources

Modificar el archivo /opt/mflow/oracle.json para poner los connection strings correspondientes al ambiente