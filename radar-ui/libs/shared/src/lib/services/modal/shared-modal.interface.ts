export interface SharedModalParams {
    title?: string | null;
    content: string;
    confirmText: string;
    cancelText: string;
    confirmHandler: () => void;
    cancelHandler?: () => void;
}
