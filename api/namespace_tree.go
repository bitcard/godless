package api

import "github.com/johnny-morrice/godless/crdt"

type NamespaceTree interface {
	JoinTable(crdt.TableName, crdt.Table) error
	LoadTraverse(NamespaceTreeTableReader) error
}

type KvNamespaceTree interface {
	KvNamespace
	NamespaceTree
}

type NamespaceTreeReader interface {
	ReadNamespace(crdt.Namespace) (bool, error)
}

type TableHinter interface {
	ReadsTables() []crdt.TableName
}

type NamespaceTreeTableReader interface {
	TableHinter
	NamespaceTreeReader
}

func AddTableHints(tables []crdt.TableName, ntr NamespaceTreeReader) NamespaceTreeTableReader {
	return tableHinterWrapper{
		hints:  tables,
		reader: ntr,
	}
}

type tableHinterWrapper struct {
	reader NamespaceTreeReader
	hints  []crdt.TableName
}

func (thw tableHinterWrapper) ReadsTables() []crdt.TableName {
	return thw.hints
}

func (thw tableHinterWrapper) ReadNamespace(ns crdt.Namespace) (bool, error) {
	return thw.reader.ReadNamespace(ns)
}

// NamespaceTreeReader functions return true when they have finished reading
// the tree.
type NamespaceTreeLambda func(ns crdt.Namespace) (bool, error)

func (ntl NamespaceTreeLambda) ReadNamespace(ns crdt.Namespace) (bool, error) {
	return ntl(ns)
}
