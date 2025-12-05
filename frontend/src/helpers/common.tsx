import type { DomainType } from "../types";

export const getHeaders = (domainType: DomainType) => {
    switch (domainType) {
      case 'onapi':
        return ['Expediente', 'Tipo', 'Subtipo', 'Texto', 'Titular'];
      case 'dgii':
        return ['RNC', 'Razón Social', 'Nombre Comercial', 'Estado', ];
      case 'scj':
        return ['Expediente', 'Sentencia', 'Tribunal', 'Materia', 'URL', 'Fecha Fallo'];
      case 'dgii':
        return ['RNC', 'Razón Social', 'Nombre Comercial', 'Estado', 'Facturador Electrónico', 'Licencia Comercial'];
      case 'pgr':
        return ['Título', 'URL'];
      case 'docking': case 'social_media': case 'file_type': case 'x_social_media': case 'google_docking':
        return ['Título', 'Descripción', 'URL'];
      default:
        return [];
    }
  };      