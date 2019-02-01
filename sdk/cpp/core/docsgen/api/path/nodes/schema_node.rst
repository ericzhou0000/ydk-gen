.. _ref-schemanode:

SchemaNode
==========


.. cpp:class:: ydk::path::SchemaNode

    Represents a Node in the SchemaTree.

    .. cpp:function:: ~SchemaNode()

        The destructor.

        .. note::

            A :cpp:class:`SchemaNode<SchemaNode>` represents a containment hierarchy. So invocation of the destructor will lead to the children of this node being destroyed.

    .. cpp:function:: virtual std::string get_path() const

        Get the path expression representing this Node in in the NodeTree.

        :return: ``std::string`` representing the path to this Node.

    .. cpp:function:: virtual std::vector<SchemaNode*> find(const std::string& path) const

        Finds descendant nodes that match the given xpath expression.

        This API finds descendant nodes in the :cpp:class:`SchemaNode<SchemaNode>` tree that satisfy the given path expression. See :ref:`how to path<ref-howtopath>`.

        :param path: The path expression.
        :return: Vector of :cpp:class:`SchemaNode<SchemaNode>` that satisfies the criterion.
        :raises: :cpp:class:`YPathError<YPathError>` if the path expression in invalid, see error code for details.
        :raises: :cpp:class:`YInvalidArgumentError<YInvalidArgumentError>` if the argument is invalid.

    .. cpp:function:: virtual const SchemaNode* get_parent() const noexcept

        Get the Parent Node of this SchemaNode in the tree.

        :return: ``SchemaNode*`` the children of this node.

    .. cpp:function:: virtual std::vector<std::unique_ptr<SchemaNode>> & get_children() const

        Returns the children of this :cpp:class:`SchemaNode <SchemaNode>` in the NodeTree.

        :return: Child schema nodes.

    .. cpp:function:: virtual const SchemaNode& get_root() const noexcept

        Get the root of NodeTree this node is part of

        :return: The pointer to the root.

    .. cpp:function:: virtual Statement get_statement() const

        Returns the YANG statement associated with this :cpp:class:`SchemaNode<SchemaNode>`

        :return: The pointer to the yang statement for this :cpp:class:`SchemaNode<SchemaNode>`

    .. cpp:function:: virtual std::vector<Statement> get_keys() const

        Returns vector of YANG statement corresponding the the keys.

        :return: Vector of :cpp:class:`Statement` that represent keys.
