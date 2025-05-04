import type { SyllabusUploadData } from '~/types/types'
import type { UploadSyllabusResult } from '~/server/api/syllabi/upload'
import { calculateChecksum } from '~/lib/utils'

export const useSyllabusUpload = () => {
  const uploadSyllabus = async (
    data: SyllabusUploadData,
    file: File
  ): Promise<{ success: boolean; errorText?: string }> => {
    try {
      const checksum = await calculateChecksum(file)
      const response = await $fetch<UploadSyllabusResult>('/api/syllabi/upload', {
        method: 'POST',
        body: {
          ...data,
          checksum,
          contentType: file.type,
          fileName: file.name,
          fileSize: file.size,
        }
      })

      // if (!response.success) {
      //   return {
      //     success: false,
      //     errorText: response.errorText,
      //   }
      // }

      const { presignedUrl } = response.data

      await $fetch(presignedUrl, {
        method: 'PUT',
        body: file,
        headers: {
          'Content-Type': file.type,
        },
      });

      return {
        success: true,
      }

    } catch (err: any) {
      console.log(err)
      return {
        success: false,
        errorText: "Unexpected error occurred. Please try again.",
      }
    }
  }

  return {
    uploadSyllabus,
  }
}
