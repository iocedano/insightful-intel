import type { DomainType } from "../types";

export const getHeaders = (domainType: DomainType) => {
    switch (domainType) {
      case 'onapi':
        return ['Expediente', 'Tipo', 'Subtipo', 'Texto', 'Titular'];
      case 'scj':
        return ['Expediente', 'Sentencia', 'Tribunal', 'Materia', 'URL', 'Fecha Fallo'];
      case 'pgr':
        return ['Título', 'URL'];
      case 'docking': case 'social_media': case 'file_type': case 'x_social_media':
        return ['Título', 'Descripción', 'Relevancia', 'URL'];
      default:
        return [];
    }
  };      