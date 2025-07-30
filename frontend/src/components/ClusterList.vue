<template>
  <div class="card">
    <h2>Clusters</h2>
    <ul v-if="clusters.length > 0">
      <li v-for="c in clusters" :key="c.name">
        <span>
          <b>{{ c.name }}</b> |
          User: {{ c.info.user }} |
          Databases: {{ c.info.databases?.join(', ') }} |
          Storage: {{ c.info.storage }} |
          NodePort: {{ c.info.nodePort || 'Pending' }}
        </span>
        <button class="get" @click="emitGet(c.name)">GET</button>
        <button @click="confirmDelete(c.name)" :disabled="loading === c.name">
          {{ loading === c.name ? 'Deleting...' : 'Delete' }}
        </button>
      </li>
    </ul>
    <p v-else style="color:#888">No clusters found.</p>
    <button class="refresh" @click="fetchClusters">Refresh</button>
    <p v-if="message" :class="messageType">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref, defineExpose } from 'vue'
const clusters = ref([]) // [{ name, info }]
const loading = ref('')
const message = ref('')
const messageType = ref('')

defineExpose({ fetchClusters })

const emit = defineEmits(['deleted', 'get'])

async function fetchClusters() {
  loading.value = ''
  message.value = ''
  try {
    const res = await fetch('http://168.119.243.127:30002/list')
    if (!res.ok) throw new Error('Error listing clusters')
    const names = await res.json()
    clusters.value = await Promise.all(names.map(async name => {
      try {
        const r = await fetch(`http://168.119.243.127:30002/get?name=${encodeURIComponent(name)}`)
        const info = r.ok ? await r.json() : {}
        return { name, info }
      } catch {
        return { name, info: {} }
      }
    }))
  } catch (err) {
    message.value = err.message
    messageType.value = 'error'
  }
}
fetchClusters()

function emitGet(name) {
  emit('get', name)
}

async function confirmDelete(name) {
  const user = window.prompt(`Enter user for cluster "${name}":`)
  if (!user) return
  const password = window.prompt(`Enter password for user "${user}":`)
  if (!password) return
  if (!confirm(`Delete cluster "${name}" as user "${user}"? This cannot be undone.`)) return
  loading.value = name
  message.value = ''
  try {
    // If your backend does NOT use user/password, just send name as before.
    // To actually check credentials, your backend must require user/password.
    // Here we send as query params (change to body if you update backend).
    const res = await fetch(
      `http://168.119.243.127:30002/delete?name=${encodeURIComponent(name)}&user=${encodeURIComponent(user)}&password=${encodeURIComponent(password)}`,
      { method: 'DELETE' }
    )
    if (!res.ok) throw new Error(await res.text())
    clusters.value = clusters.value.filter(c => c.name !== name)
    message.value = 'Cluster deleted.'
    messageType.value = 'success'
    emit('deleted')
  } catch (err) {
    message.value = err.message
    messageType.value = 'error'
  } finally {
    loading.value = ''
    setTimeout(() => (message.value = ''), 2000)
  }
}
</script>

<style scoped>
.card {
  margin-bottom: 2rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 1px 6px #0001;
  padding: 1.5rem;
}
ul {
  list-style: none;
  padding: 0;
  margin: 0 0 1rem 0;
}
li {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.4rem 0;
  border-bottom: 1px solid #f1f1f1;
  gap: 0.5rem;
}
li:last-child {
  border-bottom: none;
}
button {
  background: #e74c3c;
  color: #fff;
  border: none;
  border-radius: 5px;
  padding: 0.2rem 0.8rem;
  cursor: pointer;
  transition: background 0.2s;
}
button:disabled {
  background: #f1adad;
}
.get {
  background: #4f8cfb;
  margin-right: 0.5rem;
}
.refresh {
  margin-top: 1rem;
  background: #4f8cfb;
}
.success {
  color: #258b40;
  margin-top: 0.5rem;
}
.error {
  color: #e25241;
  margin-top: 0.5rem;
}
</style>
