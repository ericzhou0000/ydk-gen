#  ----------------------------------------------------------------
# Copyright 2016 Cisco Systems
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# ------------------------------------------------------------------

"""
source_printer.py

 prints C++ classes

"""

class ClassGetChildrenPrinter(object):
    def __init__(self, ctx):
        self.ctx = ctx

    def print_class_get_children(self, clazz, children):
        self.ctx.writeln('std::map<std::string, std::shared_ptr<ydk::Entity>> %s::get_children() const' % clazz.qualified_cpp_name())
        self.ctx.writeln('{')
        self.ctx.lvl_inc()
        self.ctx.writeln('std::map<std::string, std::shared_ptr<ydk::Entity>> _children{};')
        self.ctx.writeln('char count_=0;')
        for child in children:
            self._print_class_get_child(child)
        self.ctx.writeln('return _children;')
        self.ctx.lvl_dec()
        self.ctx.writeln('}')
        self.ctx.bline()

    def _print_class_get_child(self, child):
        if child.is_many:
            self._print_class_get_child_many(child)
        else:
            self._print_class_get_child_unique(child)
        self.ctx.bline()

    def _print_class_get_child_many(self, child):
        self.ctx.writeln('count_ = 0;')
        self.ctx.writeln('for (auto ent_ : %s.entities())' % child.name)
        self.ctx.writeln('{')
        self.ctx.lvl_inc()
        self.ctx.writeln('if(_children.find(ent_->get_segment_path()) == _children.end())')
        self.ctx.lvl_inc()
        self.ctx.writeln('_children[ent_->get_segment_path()] = ent_;')
        self.ctx.lvl_dec()
        self.ctx.writeln('else')
        self.ctx.lvl_inc()
        self.ctx.writeln('_children[ent_->get_segment_path()+count_++] = ent_;')
        self.ctx.lvl_dec()
        self.ctx.lvl_dec()
        self.ctx.writeln('}')

    def _print_class_get_child_unique(self, child):
        self.ctx.writeln('if(%s != nullptr)' % child.name)
        self.ctx.writeln('{')
        self.ctx.lvl_inc()
        path = ''
        if child.stmt.i_module.arg != child.owner.stmt.i_module.arg:
            path += child.stmt.i_module.arg
            path += ':'
        path += child.stmt.arg
        self.ctx.writeln('_children["%s"] = %s;' % (path, child.name))
        self.ctx.lvl_dec()
        self.ctx.writeln('}')
