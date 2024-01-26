import path from "path"
import react from "@vitejs/plugin-react"
import { defineConfig } from "vite"

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server:{
    proxy:{
      "/api":{
        target:"http://194.36.24.245:3100",
        changeOrigin:true,
        rewrite:(path)=>path.replace(/^\/api/,"")
      }
    }
  },
})
