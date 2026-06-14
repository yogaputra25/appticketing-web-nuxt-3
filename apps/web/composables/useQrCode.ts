import QRCode from 'qrcode'

export function useQrCode() {
  const config = useRuntimeConfig()

  function getVerifyUrl(ticketCode: string): string {
    const base = config.public.apiBase || 'http://localhost:8080'
    return `${base}/tickets/v/${ticketCode}`
  }

  async function generateDataUrl(ticketCode: string): Promise<string> {
    const url = getVerifyUrl(ticketCode)
    return QRCode.toDataURL(url, {
      width: 300,
      margin: 2,
      color: {
        dark: '#1e293b',
        light: '#ffffff',
      },
    })
  }

  return {
    getVerifyUrl,
    generateDataUrl,
  }
}
