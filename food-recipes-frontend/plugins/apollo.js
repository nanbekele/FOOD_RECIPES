import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context'
import { DefaultApolloClient } from '@vue/apollo-composable'
import { useRuntimeConfig } from '#app'

export default defineNuxtPlugin((nuxtApp) => {
  const runtime = useRuntimeConfig()
  const hasuraUri = runtime.public?.NUXT_PUBLIC_HASURA_ENDPOINT || 'http://localhost:8082/v1/graphql'
  const httpLink = createHttpLink({ uri: hasuraUri })

  const authLink = setContext((_, { headers }) => {
    // Only access localStorage on the client
    let token = ''
    if (typeof window !== 'undefined') {
      token = localStorage.getItem('token') || ''
    }
    const nextHeaders = { ...headers }
    // Rely on JWT claims for role; do not override x-hasura-role from client
    if (token) nextHeaders.Authorization = `Bearer ${token}`
    return { headers: nextHeaders }
  })

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
  })

  nuxtApp.vueApp.provide(DefaultApolloClient, apolloClient)
})
