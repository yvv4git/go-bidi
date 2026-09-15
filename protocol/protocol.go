// Package protocol implements the WebDriver BiDi wire protocol: message
// types and their JSON encoding, as exchanged between a client and a
// BiDi-enabled browser.
//
// Commands carry an id, a method and params. Responses match the command id
// and carry either a result or an error. Events carry a method and params
// without an id. Framed as JSON, and tagged with a "type" field on the
// browser side:
//
//	{ "id": 1, "method": "session.new", "params": { ... } }
//	{ "type": "success", "id": 1, "result": { ... } }
//	{ "type": "error", "id": 1, "error": "...", "message": "..." }
//	{ "type": "event", "method": "log.entryAdded", "params": { ... } }
package protocol
