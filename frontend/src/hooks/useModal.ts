import { useState, useCallback } from "react";
interface UseModalResult<T> {
    isOpen: boolean;
    data: T | null;
    open: (data: T) => void;
    close: () => void;
}
export const useModal = <T>(): UseModalResult<T> => {
    const [isOpen, setIsOpen] = useState(false);
    const [data, setData] = useState<T | null>(null);
    const open = useCallback((modalData: T) => {
        setData(modalData);
        setIsOpen(true);
    }, []);
    const close = useCallback(() => {
        setIsOpen(false);
        setData(null);
    }, []);
    return { isOpen, data, open, close };
};
