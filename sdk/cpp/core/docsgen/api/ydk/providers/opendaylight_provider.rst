OpenDaylightServiceProvider
============================


.. cpp:class:: ydk::OpenDaylightServiceProvider

  A service provider to be used to communicate with an  `OpenDaylight <https://opendaylight.org>`_ instance.

    .. cpp:function::
        OpenDaylightServiceProvider(path::Repository & repo, \
                                    const std::string & address, \
                                    const std::string & username, \
                                    const std::string & password, \
                                    int port = 8181, \
                                    EncodingFormat encoding = EncodingFormat::JSON, \
                                    Protocol protocol = Protocol::restconf)

        Constructs an instance of the ``OpenDaylightServiceProvider`` to connect to a OpenDaylight instance

        :param repository: Reference to an instance of :cpp:class:`path::Repository<ydk::path::Repository>`
        :param address: IP address of the ODL instance
        :param username: Username to log in to the instance
        :param password: Password to log in to the instance
        :param port: Device port used to access the ODL instance. Default value is 8181
        :param encoding: Type of encoding to be used for the payload. Default is :cpp:enumerator:`JSON<EncodingFormat::JSON>`
        :param protocol: Type of OpenDaylight northbound protocol. Currently, only restconf is supported and is the default value

    .. cpp:function:: ydk::ServiceProvider & get_node_provider(const std::string & node_id)

        Returns the :cpp:class:`ServiceProvider<ydk::ServiceProvider>` instance corresponding to the device being controlled by the OpenDaylight instance, indicated by "node_id".

        :param node_id: The name of the device being controlled by the OpenDaylight instance.
        :return: Reference to the :cpp:class:`ServiceProvider<ydk::ServiceProvider>` or raises :cpp:class:`YServiceProviderError<YServiceProviderError>` if one could not be found.

    .. cpp:function:: const std::vector<std::string> & get_node_ids()

         Returns a list of node ID's of the devices being controlled by this OpenDaylight instance.

        :return: List of node ID's of the devices being controlled by this OpenDaylight instance.

    .. cpp:function:: ~OpenDaylightServiceProvider()
