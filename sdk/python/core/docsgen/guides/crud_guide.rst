.. _netconf-operations:

How do I create, update, read and delete?
=========================================

.. contents:: Table of Contents

This document contains some examples of creating, reading and deleting yang data using YDK. To perform these operations, the :py:class:`CRUDService<ydk.services.CRUDService>` is used. Also, in these examples, :py:class:`YFilter<ydk.filters.YFilter>` is used to mark parts of the data for particular operations.

Creating a configuration with a list and a presence class
---------------------------------------------------------

To configure a rule in the SNMP trap correlator, the below approach can be used.

Note that the :py:class:`Cisco_IOS_XR_snmp_agent_cfg.Snmp.Correlator.Rules.Rule<ydk.models.cisco_ios_xr.Cisco_IOS_XR_snmp_agent_cfg.Snmp.Correlator.Rules.Rule>` class is a :py:class:`YList<ydk.types.YList>`. So it needs to be instantiated and appended to its parent.

Also, the attribute ``non_stateful`` is an instance of a :ref:`presence class<presence-class>`, which is set to ``None`` by default. So it needs to be assigned to a new instance of its class.

.. code-block:: python
    :linenos:

    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_snmp_agent_cfg

    # Create the top-level container
    snmp = Cisco_IOS_XR_snmp_agent_cfg.Snmp()

    # Create the list instance
    rule = Cisco_IOS_XR_snmp_agent_cfg.Snmp.Correlator.Rules.Rule()
    rule.name = 'abc'

    # Instantiate and assign the presence class
    rule.non_stateful = Cisco_IOS_XR_snmp_agent_cfg.Snmp.Correlator.Rules.Rule.NonStateful()

    rule.non_stateful.timeout = 3

    # Append the list instance to its parent
    snmp.correlator.rules.rule.append(rule)

    # Call the CRUD create on the top-level snmp object
    # (assuming you have already instantiated the service and provider)
    result = crud.create(provider, snmp)

Creating and replacing a configuration
--------------------------------------

First, let us create a configuration for the :py:class:`openconfig_bgp.Bgp.Global_.Config<ydk.models.openconfig.openconfig_bgp.Bgp.Global_.Config>` class. Here, we set the leaf ``as_``, which represents the autonomous system number, to ``65001`` and the leaf ``router_id`` to ``'10.0.0.1'``.

.. code-block:: python
    :linenos:

    from ydk.types import YFilter
    from ydk.models.openconfig import openconfig_bgp

    # First, create the top-level Bgp() object
    bgp = openconfig_bgp.Bgp()

    # Populate the values for the global config object
    bgp.global_.config.as_ = 65001
    bgp.global_.config.router_id = '10.0.0.1'

    # Call the CRUD create on the top-level bgp object
    # (assuming you have already instantiated the service and provider)
    result = crud.create(provider, bgp)

Now, let us replace the above configuration with a new configuration for the :py:class:`openconfig_bgp.Bgp.Global_.Config<ydk.models.openconfig.openconfig_bgp.Bgp.Global_.Config>` class using the below code.

.. code-block:: python
    :linenos:

    from ydk.types import YFilter
    from ydk.models.openconfig import openconfig_bgp

    # First, create the top-level Bgp() object
    bgp = openconfig_bgp.Bgp()

    # Set the yfilter attribute of the config object to YFilter.replace
    bgp.global_.config.yfilter = YFilter.replace

    # Populate the new values for the global config object
    bgp.global_.config.as_ = 65023
    bgp.global_.config.router_id = '25.3.55.12'

    # Call the CRUD update on the top-level bgp object
    # (assuming you have already instantiated the service and provider)
    result = crud.update(provider, bgp)


Creating and reading a list
---------------------------

For example, to read the instances of a deeply nested :py:class:`YList<ydk.types.YList>` called :py:class:`Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af.Safs.Saf.IpRibRouteTableNames.IpRibRouteTableName.Routes.Route<ydk.models.cisco_ios_xr.Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af.Safs.Saf.IpRibRouteTableNames.IpRibRouteTableName.Routes.Route>`  in the ``Cisco_IOS_XR_ip_rib_ipv4_oper`` module using YDK's :py:class:`CRUDService<ydk.services.CRUDService>`, the below approach can be used.

.. code-block:: python
    :linenos:

    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_ip_rib_ipv4_oper
    from ydk.filters import YFilter

    # First create the top-level Rib() object
    rib = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib()

    # Then create the list instance Vrf()
    vrf = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf()
    vrf.vrf_name='default'

    # Then create the child list element Af() and the rest of the nested list instances
    af = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af()
    af.af_name = 'IPv4'

    saf = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af.Safs.Saf()
    saf.saf_name='Unicast'

    table_name = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af.Safs.Saf.IpRibRouteTableNames.IpRibRouteTableName()
    table_name.route_table_name = 'default'

    # Create the final list instance Route()
    route = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib.Vrfs.Vrf.Afs.Af.Safs.Saf.IpRibRouteTableNames.IpRibRouteTableName.Routes.Route()
    route.yfilter = YFilter.read # set the yfilter attribute for route to YFilter.read

    # Append each of the list instances to their respective parents
    table_name.routes.route.append(route)
    saf.ip_rib_route_table_names.ip_rib_route_table_name.append(table_name)
    af.safs.saf.append(saf)
    vrf.afs.af.append(af)
    rib.vrfs.vrf.append(vrf)

    # Call the CRUD read on the top-level rib object
    # (assuming you have already instantiated the service and provider)
    rib_oper = crud.read(provider, rib)

    vrf_list = rib_oper.vrfs.vrf
    vrf_default = vrf_list["default"]  # or  vrf_default = vrf_list.get("default")

Read all VRF configuration:

.. code-block:: python
    :linenos:

    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_ip_rib_ipv4_oper
    from ydk.filters import YFilter

    # First create the top-level Rib() object
    rib = Cisco_IOS_XR_ip_rib_ipv4_oper.Rib()

    # Call the CRUD read on the top-level rib object
    # (assuming you have already instantiated the service and provider)
    rib_oper = crud.read(provider, rib)

    # Access all VRFs in the list
    for vrf in rib_oper.vrfs.vrf:
        print(vrf.vrf_name)

    # Get list of VRF names
    all_vrf_names = rib_oper.vrfs.vrf.keys()

    # Iterate over VRF names
    for vrf_name in all_vrf_names:
        vrf = rib_oper.vrfs.vrf[vrf_name]
        for af in vrf.afs.af:
            print("VRF: %s, AF: %s", vrf_name, af)

    # Access specific VRF, when name is known
    vrf = rib_oper.vrfs.vrf["default"]
    if vrf is not None:
        for af in vrf.afs.af:
            print("VRF: %s, AF: %s", "default", af)

Reading a leaf
--------------

For example, to read a :py:class:`YLeaf<ydk.types.YLeaf>` called ``running`` in the :py:class:`Cisco_IOS_XR_clns_isis_cfg.Isis.Instances.Instance<ydk.models.cisco_ios_xr.Cisco_IOS_XR_clns_isis_cfg.Isis.Instances.Instance>` class in the ``Cisco_IOS_XR_clns_isis_cfg`` module using YDK's :py:class:`CRUDService<ydk.services.CRUDService>`, the below approach can be used.

.. code-block:: python
    :linenos:


    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_clns_isis_cfg
    from ydk.types import Empty
    from ydk.filters import YFilter

    # First create the top-level Isis() object
    isis = Cisco_IOS_XR_clns_isis_cfg.Isis()

    # Create the list instance
    ins = Cisco_IOS_XR_clns_isis_cfg.Isis.Instances.Instance()
    ins.instance_name = 'default'

    # Set the yfilter attribute of the leaf called 'running' to YFilter.read
    ins.running = YFilter.read

    # Append the instance to the parent
    isis.instances.instance.append(ins)

    # Call the CRUD read on the top-level isis object
    # (assuming you have already instantiated the service and provider)
    result = crud.read(provider, isis)


Deleting a list
---------------

For example, to delete a :py:class:`YList<ydk.types.YList>` called :py:class:`Instance<ydk.models.cisco_ios_xr.Cisco_IOS_XR_clns_isis_cfg.Isis>` in the ``Cisco_IOS_XR_clns_isis_cfg`` module using YDK's :py:class:`CRUDService<ydk.services.CRUDService>`, the below approach can be used.

.. code-block:: python
    :linenos:


    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_clns_isis_cfg
    from ydk.types import Empty
    from ydk.filters import YFilter

    # First create the top-level Isis() object
    isis = Cisco_IOS_XR_clns_isis_cfg.Isis()

    # Create the list instance
    ins = Cisco_IOS_XR_clns_isis_cfg.Isis.Instances.Instance()
    ins.instance_name = 'xyz'

    # Set the yfilter attribute of the leaf called 'running' to YFilter.delete
    ins.yfilter = YFilter.delete

    # Append the instance to the parent
    isis.instances.instance.append(ins)

    # Call the CRUD update on the top-level isis object
    # (assuming you have already instantiated the service and provider)
    result = crud.update(provider, isis)



Deleting a leaf
---------------

For example, to delete a :py:class:`YLeaf<ydk.types.YLeaf>` called ``timer`` of type ``int`` in the :py:class:`Cdp<ydk.models.cisco_ios_xr.Cisco_IOS_XR_cdp_cfg.Cdp>` class in the ``Cisco_IOS_XR_cdp_cfg`` module using YDK's :py:class:`CRUDService<ydk.services.CRUDService>`, the below approach can be used.

.. code-block:: python
    :linenos:


    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_cdp_cfg
    from ydk.filters import YFilter

    # First create the top-level Cdp() object
    cdp = Cisco_IOS_XR_cdp_cfg.Cdp()

    # Set a dummy value to the leaf
    cdp.timer = 5
    # Assign YFilter.delete to the 'timer' leaf
    cdp.timer = YFilter.delete

    # Call the CRUD update on the top-level cdp object
    # (assuming you have already instantiated the service and provider)
    result = crud.update(provider, cdp)

For example, to delete a :py:class:`YLeaf<ydk.types.YLeaf>` called ``running`` of type ``Empty`` in the :py:class:`Instance<ydk.models.cisco_ios_xr.Cisco_IOS_XR_clns_isis_cfg.Isis>` class in the ``Cisco_IOS_XR_clns_isis_cfg`` module using YDK's :py:class:`CRUDService<ydk.services.CRUDService>`, the below approach can be used.

.. code-block:: python
    :linenos:


    from ydk.models.cisco_ios_xr import Cisco_IOS_XR_clns_isis_cfg
    from ydk.types import Empty
    from ydk.filters import YFilter

    # First create the top-level Isis() object
    isis = Cisco_IOS_XR_clns_isis_cfg.Isis()

    # Create the list instance
    ins = Cisco_IOS_XR_clns_isis_cfg.Isis.Instances.Instance()
    ins.instance_name = 'default'

    # Assign YFilter.delete to the 'running' leaf
    ins.running = YFilter.delete

    # Append the instance to the parent
    isis.instances.instance.append(ins)

    # Call the CRUD update on the top-level isis object
    # (assuming you have already instantiated the service and provider)
    result = crud.update(provider, isis)

    
Applying CRUD to multiple entities
--------------------------------------

You can apply CRUD operations on multiple entities in one Crud-service call. For example, you want to 'read' BGP and Interfaces configuration together.

.. code-block:: python
    :linenos:

    from ydk.types import Filter, Config
    from ydk.models.openconfig import openconfig_bgp, openconfig_interfaces

    # First, create the top-level Bgp and Interface objects
    int_filter = openconfig_interfaces.Interfaces()
    bgp_filter = openconfig_bgp.Bgp()
    
    # Create read filter
    read_filter = Filter(int_filter, bgp_filter)

    # Call the CRUD read-config to get configuration of entities
    result = crud.read_config(provider, read_filter)
    
    # Access read results from returned Config collection
    int_config = result[int_filter]
    bgp_config = result[bgp_filter]
    
    # Or print all configuration in XML format
    codec_service = CodecService()
    codec_provider = CodecServiceProvider()
    codec_provider.encoding = EncodingFormat.XML
    for entity in result:
        xml_encode = codec_service.encode(codec_provider, entity)
        print(xml_encode)
