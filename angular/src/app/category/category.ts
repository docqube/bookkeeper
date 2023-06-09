export interface Category {
  id: number;
  name: string;
  description?: string;
  color?: string;
}

export interface CategoryRule {
  id: number;
  categoryID: number;
  description: string;
  regex: string;
  mappingField: string;
}

export const CategoryColors = [
  "#4f46e5",
  "#9333ea",
  "#c026d3",
  "#ec4899",
  "#e11d48",
  "#ea580c",
  "#d97706",
  "#eab308",
];
