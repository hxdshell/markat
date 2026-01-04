const API_URL = 'http://localhost:3000/api'
export const client = {
  get: async (path: string): Promise<ApiResponseType> => {
    const res = await fetch(`${API_URL}${path}`)
    if (!res.ok) {
      const contentType = res.headers.get('Content-Type')
      if (contentType !== 'application/json') {
        const cause = await res.text()
        return {
          status: res.status,
          message: 'get request failed',
          data: null,
          cause: cause,
        }
      }
    }
    const resData = await res.json()
    return { ...resData, status: res.status }
  },
  put: async (path: string, reqData: any): Promise<ApiResponseType> => {
    const res = await fetch(`${API_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(reqData),
    })
    if (!res.ok) {
      const contentType = res.headers.get('Content-Type')
      if (contentType !== 'application/json') {
        const cause = await res.text()
        return {
          status: res.status,
          message: 'get request failed',
          data: null,
          cause: cause,
        }
      }
    }
    const resData = await res.json()
    return { ...resData, status: res.status }
  },
}
