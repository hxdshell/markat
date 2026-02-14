export const API_URL = 'http://localhost:3000/api'
export const client = {
  get: async (path: string): Promise<ApiResponseType> => {
    const res = await fetch(`${API_URL}${path}`)
    const contentType = res.headers.get('Content-Type')
    if (!res.ok) {
      if (contentType !== 'application/json') {
        const cause = await res.text()
        return {
          status: res.status,
          message: 'get request failed',
          data: null,
          cause: cause,
          contentType: contentType ?? undefined,
        }
      }
    }
    if (contentType === 'application/json') {
      const resData = await res.json()
      return {
        ...resData,
        status: res.status,
        contentType: contentType ?? undefined,
      }
    } else {
      const resData = await res.text()
      return {
        message: 'success',
        data: resData,
        status: res.status,
        contentType: contentType ?? undefined,
      }
    }
  },
  put: async (path: string, reqData: any): Promise<ApiResponseType> => {
    const res = await fetch(`${API_URL}${path}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(reqData),
    })
    const contentType = res.headers.get('Content-Type')
    if (!res.ok) {
      if (contentType !== 'application/json') {
        const cause = await res.text()
        return {
          status: res.status,
          message: 'get request failed',
          data: null,
          cause: cause,
          contentType: contentType ?? undefined,
        }
      }
    }
    if (contentType === 'application/json') {
      const resData = await res.json()
      return {
        ...resData,
        status: res.status,
        contentType: contentType ?? undefined,
      }
    } else {
      const resData = await res.text()
      return {
        message: 'success',
        data: resData,
        status: res.status,
        contentType: contentType ?? undefined,
      }
    }
  },
}
