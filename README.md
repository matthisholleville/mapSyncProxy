<a href="https://www.haproxy.com">
    <img src="https://upload.wikimedia.org/wikipedia/commons/a/ab/Haproxy-logo.png" alt="Pritunl logo" title="Pritunl" align="right" height="100" />
</a>
<a href="https://cloud.google.com/storage">
    <img src="https://imgs.search.brave.com/5lxcXp7DQkSKquKVb6CQapUrgQTRsibWDzbcaLBqfi0/rs:fit:860:0:0/g:ce/aHR0cHM6Ly9zdGF0/aWMtMDAuaWNvbmR1/Y2suY29tL2Fzc2V0/cy4wMC9jbG91ZC1z/dG9yYWdlLWljb24t/MjU2eDIwNC1kb3Z3/cGp5eC5wbmc" alt="GCS logo" title="GCS" align="right" height="100" />
</a>

# Simple proxy to synchronize a GCS file with an HAProxy .map file

This small Golang server exposes a route that allows synchronizing a GCS (Google Cloud Storage) file with an HAProxy .map file using the Dataplane API.

## Use cases

This proxy can be used for various purposes:

1. Implement distributed rate limiting across multiple HAProxy instances without the need for an enterprise license.
2. Apply access rules (allowlisting/denylisting) across all machines in one or more environments.
3. Ensure the reliable updating of Map files.

## Demo
