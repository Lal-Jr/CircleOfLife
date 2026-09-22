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
        </header>
    );
}
