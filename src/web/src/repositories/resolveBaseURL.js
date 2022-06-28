function resolveBaseURL(serviceDomain, servicePort, servicePath) {
  if (process.env.NODE_ENV === 'production') {
    return ""
  }

  servicePort = servicePort || ''
  servicePath = servicePath || '/'
  return `${serviceDomain}${servicePort ? `:${servicePort}` : ''}${servicePath}`
}

export default resolveBaseURL