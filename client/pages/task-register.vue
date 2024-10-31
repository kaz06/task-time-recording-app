<template>
    <v-app>
      <v-main>
        <v-container>
          <!-- Task Registration Form -->
          <v-form @submit.prevent="submitTime">
            <v-text-field
              v-model="title"
              label="Task Title"
              placeholder="Enter the task title"
              required
            ></v-text-field>
  
            <v-date-picker
              v-model="date"
              label="Task Completion Date"
              required
            ></v-date-picker>
  
            <v-text-field
              v-model="time"
              label="Task Time"
              type="time"
              required
            ></v-text-field>
  
            <v-text-field
              v-model="tags"
              label="Tags"
              placeholder="Enter tags separated by commas"
              required
            ></v-text-field>
  
            <v-btn type="submit">Submit</v-btn>
          </v-form>
  
          <!-- Show Chart Button -->
          <v-btn @click="chartDialog = true">Show Chart</v-btn>
  
          <!-- Chart Display Modal -->
          <v-dialog v-model="chartDialog" max-width="800">
            <v-card>
              <v-card-title>Chart Display</v-card-title>
              <v-card-text>
                <v-form @submit.prevent="getTimeByTag">
                  <v-text-field
                    v-model="startDate"
                    label="Start Date"
                    type="date"
                    required
                  ></v-text-field>
  
                  <v-text-field
                    v-model="endDate"
                    label="End Date"
                    type="date"
                    required
                  ></v-text-field>
  
                  <v-btn type="submit">Get Time by Tag</v-btn>
                </v-form>
  
                <div v-if="Object.keys(timeByTag).length > 0">
                  <Pie :data="chartData" :options="chartOptions" />
                </div>
              </v-card-text>
            </v-card>
          </v-dialog>
        </v-container>
      </v-main>
    </v-app>
  </template>
  
  <script setup lang="ts">
  import { ref, computed } from 'vue'
  import { useRuntimeConfig, useFetch } from '#app'
  import { Pie } from 'vue-chartjs'
  import { Chart as ChartJS, Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale } from 'chart.js'
  import { useAuth } from '@/composables/auth'


  ChartJS.register(Title, Tooltip, Legend, ArcElement, CategoryScale, LinearScale)
  
  const { getToken } = useAuth()
  const config = useRuntimeConfig()
  
  // Form Data
  const title = ref('')
  const time = ref('00:01')
  const tags = ref('')
  const date = ref<Date | null>(null)
  
  // For Chart Display
  const startDate = ref('')
  const endDate = ref('')
  const timeByTag = ref<{ [key: string]: number }>({})
  const chartDialog = ref(false)
  
  const chartData = computed(() => ({
    labels: Object.keys(timeByTag.value),
    datasets: [
      {
        label: 'Task Time by Tag',
        data: Object.values(timeByTag.value),
        backgroundColor: [
          '#FF6384',
          '#36A2EB',
          '#FFCE56',
          '#4BC0C0',
          '#9966FF',
          '#FF9F40'
        ],
      }
    ]
  }))
  
  const chartOptions = {
    responsive: true
  }
  
  const submitTime = async () => {
    try {
      const idToken = await getToken()
      if (!idToken) {
        console.error('ID token not found.')
        return
      }


      const formattedTime = time.value.length === 5 ? `${time.value}:00` : time.value
  
      const { data, error } = await useFetch<{ time: string }>('/v1/tasks', {
        baseURL: config.public.baseURL as string,
        method: 'POST',
        body: {
          title: title.value,
          task_time: formattedTime,
          task_finish_date: date.value?.toISOString(),
          tags: tags.value.split(',').map(tag => tag.trim())
        },
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
        credentials: 'include'
      })
      if (error.value) {
        console.error('Task submission error:', error.value)
      } else {
        console.log('Task successfully submitted:', data.value)
      }
    } catch (error) {
      console.error('Task submission error:', error)
    }
  }
  
  const getTimeByTag = async () => {
    try {
      const idToken = await getToken()
      if (!idToken) {
        console.error('ID token not found.')
        return
      }
      const { data, error } = await useFetch<{ [key: string]: number }>('/v1/tasktimebytag', {
        baseURL: config.public.baseURL as string,
        method: 'GET',
        headers: {
          Authorization: `Bearer ${idToken}`,
        },
        params: {
          start_date: startDate.value,
          end_date: endDate.value
        },
        credentials: 'include'
      })
      if (error.value) {
        console.error('Error fetching time by tag:', error.value)
      } else {
        timeByTag.value = data.value || {}
  
        if (Object.keys(timeByTag.value).length === 0) {
          console.error('No data available for the specified period.')
        }
      }
    } catch (error) {
      console.error('Error fetching time by tag:', error)
    }
  }
  </script>
