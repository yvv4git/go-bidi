// Package protocol implements the WebDriver BiDi wire protocol: message
// types and their JSON-RPC 2.0 encoding, as exchanged between a client and
// a BiDi-enabled browser.
//
// Commands carry an id, a method and params; responses match the command id
// and carry either a result or an error; events carry a method and params
// without an id:
//
//	{ "id": 1, "method": "session.new", "params": { ... } }
//	{ "id": 1, "result": { ... } }
//	{ "method": "log.entryAdded", "params": { ... } }
package protocol
