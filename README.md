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


<img width="1105" height="1193" alt="Screenshot 2026-02-16 at 10 29 45 AM" src="https://github.com/user-attachments/assets/641258a7-8ec2-4447-bc85-edf6e4a313fb" />
<img width="1268" height="587" alt="Screenshot 2026-02-16 at 10 30 01 AM" src="https://github.com/user-attachments/assets/783d4491-d705-4951-8e44-af47fe24c31a" />
