// Config loader - loads configuration from YAML file based on environment
import yaml from 'js-yaml'

// Import all config files explicitly
import devConfigContent from '../../config.development.yaml?raw'
import prodConfigContent from '../../config.production.yaml?raw'
import defaultConfigContent from '../../config.yaml?raw'

// Determine which config to use based on environment
const env = import.meta.env.MODE || 'development'

// Select the appropriate config content
let yamlContent
if (env === 'production') {
  yamlContent = prodConfigContent
} else if (env === 'development') {
  yamlContent = devConfigContent
} else {
  yamlContent = defaultConfigContent
}

const parsedConfig = yaml.load(yamlContent)

const config = {
  // Current environment
  ENV: env,

  // API Base URL from YAML config
  API_BASE_URL: parsedConfig.API_BASE_URL || 'http://localhost:8081/api',

  // Upload URL - used for file uploads
  get UPLOAD_URL() {
    return `${this.API_BASE_URL}/upload`
  },

  // Get upload URL with type parameter (products, news)
  getUploadUrl(type = 'products') {
    return `${this.API_BASE_URL}/upload?type=${type}`
  }
}

export default config
