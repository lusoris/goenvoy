// Package whisparr provides clients for the Whisparr v2 and v3 (Eros) APIs.
//
// Whisparr v2 is Sonarr-based and models content as series (sites) with
// episodes (scenes). Whisparr v3, codenamed Eros, is Radarr-based and
// models content as individual movies/scenes with full performer and
// studio management.
//
// Use [Client] for Whisparr v2 instances and [ErosClient] for v3/Eros
// instances. Both wrap [arr.BaseClient] and expose typed methods for
// their respective resources.
package whisparr
