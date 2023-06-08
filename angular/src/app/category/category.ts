

export interface Category {
  id: number;
  name: string;
  description: string;
  color: string;
}

export interface CategoryRule {
  id: number;
  categoryID: number;
  description: string;
  regex: string;
  mappingField: string;
}
