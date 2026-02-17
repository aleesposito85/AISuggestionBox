**Suggestion Box web application**

- website build in simple html, js, css, (deployed with caddy server on a proxmox container)
- the backend API build in golang (deployed on a proxmox container)
- database is sql server (deployed with docker)
- N8N is used to trigger the AI model (deployed on a proxmox container)
- nodered is used to send the emails (deployed on a proxmox container)
- AI model running with ollama (deployed with docker)
 

**This AI suggestion box, creates 3 replies:**
- one for a generic reply, that you can see on the website itself if refresh the page
- one as self actions the user can take sent to the user email
- one for the leadership as a plan of actions sent to the leadership team email
