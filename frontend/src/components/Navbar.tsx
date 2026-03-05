"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { MapPin, Home, PlusCircle, User, LogOut } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { Button } from "@/components/ui/button";

export function Navbar() {
    const pathname = usePathname();
    const { logout } = useAuth();

    // Hide Navbar completely on auth pages
    if (pathname === "/login" || pathname === "/signup") {
        return null;
    }

    const navItems = [
        { name: "Feed", href: "/feed", icon: Home },
        { name: "Create Post", href: "/create", icon: PlusCircle },
        { name: "Profile", href: "/profile", icon: User },
    ];

    return (
        <header className="sticky top-0 z-50 w-full border-b bg-background/80 backdrop-blur-md">
            <div className="container mx-auto flex h-16 items-center justify-between px-4 md:px-6">
                <Link href="/" className="flex items-center gap-2 transition-opacity hover:opacity-80">
                    <div className="flex h-8 w-8 items-center justify-center rounded-full bg-primary/10 text-primary">
                        <MapPin className="h-5 w-5" />
                    </div>
                    <span className="font-bold tracking-tight text-lg">Circle</span>
                </Link>

                {/* Desktop Navigation */}
                <nav className="hidden md:flex gap-6 items-center flex-1 justify-center">
                    {navItems.map((item) => {
                        const isActive = pathname === item.href;
                        return (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`flex items-center gap-2 text-sm font-medium transition-colors hover:text-primary ${isActive ? "text-primary" : "text-muted-foreground"
                                    }`}
                            >
                                <item.icon className="h-4 w-4" />
                                {item.name}
                            </Link>
                        );
                    })}
                </nav>

                <div className="flex items-center gap-4">
                    <Button variant="ghost" size="sm" onClick={logout} className="hidden md:flex gap-2">
                        <LogOut className="h-4 w-4" />
                        Logout
                    </Button>
                </div>
            </div>

            {/* Mobile Bottom Navigation Bar */}
            <div className="md:hidden fixed bottom-0 left-0 right-0 border-t bg-background/90 backdrop-blur-md pb-safe">
                <nav className="flex justify-around items-center h-16 px-2">
                    {navItems.map((item) => {
                        const isActive = pathname === item.href;
                        return (
                            <Link
                                key={item.href}
                                href={item.href}
                                className={`flex flex-col items-center justify-center w-full h-full gap-1 transition-colors ${isActive ? "text-primary" : "text-muted-foreground hover:text-primary"
                                    }`}
                            >
                                <item.icon className="h-5 w-5" />
                                <span className="text-[10px] font-medium">{item.name}</span>
                            </Link>
                        );
                    })}
                    <button
                        onClick={logout}
                        className="flex flex-col items-center justify-center w-full h-full gap-1 text-muted-foreground hover:text-primary transition-colors"
                    >
                        <LogOut className="h-5 w-5" />
                        <span className="text-[10px] font-medium">Logout</span>
                    </button>
                </nav>
            </div>
        </header>
    );
}
