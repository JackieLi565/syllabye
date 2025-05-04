import { SyllabusUploadData } from '~/types/types';

export interface UploadSyllabusSuccess {
  success: true;
  data: {
    body: string;
    location: string;
    presignedUrl: string;
  };
}

export interface UploadSyllabusFailure {
  success: false;
  status: number;
  errorText: string;
}

export type UploadSyllabusResult = UploadSyllabusSuccess | UploadSyllabusFailure;

export default defineEventHandler(async (event) => {
  // const uploadFile = async (data: SyllabusUploadData): Promise<UploadSyllabusResult> => {
  const config = useRuntimeConfig();
  const data = await readBody(event);
  const apiUrl = `${config.public.apiUrl}/syllabi`;
  const cookie = getRequestHeader(event, "cookie");

  try {
    const response = await $fetch.raw<string>(apiUrl, {
      method: 'POST',
      headers: {
        cookie: cookie || "",
      },
      body: data,
    });

    const presignedUrl = response.headers.get('X-Presigned-Url') ?? '';
    const location = response.headers.get('Location') ?? '';
    const body = await response.text();

    return {
      success: true,
      data: {
        body,
        location,
        presignedUrl,
      },
    };
  } catch (err: any) {
    console.log(err);
    return null
  }
});
// })