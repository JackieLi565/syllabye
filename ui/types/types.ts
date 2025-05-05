export interface Session {
  userId: string;
}

export interface User {
  id: string;
  currentYear?: number | null;
  email: string;
  fullname: string;
  gender?: "male" | "female" | "other";
  nickname?: string;
  picture: string;
  programId?: string;
  instagram?: string;
  bio?: string;
  newuser: boolean;
}

export interface Program {
  faculty: string;
  id: string;
  name: string;
  uri: string;
}

export interface Syllabus {
  contentType: string;
  courseId: string;
  dateAdded: number;
  fileName: string;
  fileSize: number;
  id?: string;
  received: boolean;
  semester: string;
  userId: string;
  year: number;
}

export interface Category {
  // unique identifier
  id: string;
  // name of the category
  name: string;
}

export interface Course {
  id: string;
  categoryId: string;
  description: string;
  course: string;
  title: string;
  uri: string;
}
