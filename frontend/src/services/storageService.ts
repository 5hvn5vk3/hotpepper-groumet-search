class StorageService {
    get<T>(key: string): T | null {
        try {
            const item = sessionStorage.getItem(key);
            return item ? JSON.parse(item) : null;
        }
        catch (error) {
            console.error(`Failed to get item from storage: ${key}`, error);
            return null;
        }
    }
    set<T>(key: string, value: T): void {
        try {
            sessionStorage.setItem(key, JSON.stringify(value));
        }
        catch (error) {
            console.error(`Failed to set item in storage: ${key}`, error);
        }
    }
    remove(key: string): void {
        try {
            sessionStorage.removeItem(key);
        }
        catch (error) {
            console.error(`Failed to remove item from storage: ${key}`, error);
        }
    }
}
export const storageService = new StorageService();
