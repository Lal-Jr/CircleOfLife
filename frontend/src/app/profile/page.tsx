"use client";

import { useCurrentUser } from "@/hooks/useCurrentUser";
import { useRequireAuth } from "@/hooks/useRequireAuth";
import { useAuth } from "@/hooks/useAuth";
import { formatTimeAgo } from "@/lib/utils";
import { AlertCircle, LogOut, Mail, User as UserIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Card, CardContent } from "@/components/ui/card";

export default function ProfilePage() {
    const { isAuthChecked } = useRequireAuth();
    const { user, isLoading, error } = useCurrentUser();
    const { logout } = useAuth();

    if (!isAuthChecked) return null;

    return (
        <div className="container max-w-2xl mx-auto px-4 py-8 animate-in fade-in duration-500">
            <div className="mb-8">
                <h1 className="text-3xl font-bold tracking-tight">Profile</h1>
                <p className="text-muted-foreground mt-1">Your account details</p>
            </div>

            {isLoading ? (
                <Card className="border-muted shadow-sm">
                    <CardContent className="p-6 animate-pulse space-y-4">
                        <div className="flex items-center gap-4">
                            <div className="h-16 w-16 rounded-full bg-muted"></div>
                            <div className="space-y-2 flex-1">
                                <div className="h-4 bg-muted rounded w-1/3"></div>
                                <div className="h-3 bg-muted rounded w-1/2"></div>
                            </div>
                        </div>
                    </CardContent>
                </Card>
            ) : error || !user ? (
                <div className="flex flex-col items-center justify-center py-16 text-center bg-destructive/5 rounded-xl border border-destructive/20">
                    <AlertCircle className="h-10 w-10 text-destructive mb-4" />
                    <h3 className="text-lg font-semibold mb-2">Could not load your profile</h3>
                    <p className="text-muted-foreground text-sm mb-6">Please try again in a moment.</p>
                    <Button variant="outline" onClick={() => window.location.reload()}>Retry</Button>
                </div>
            ) : (
                <Card className="border-muted shadow-md overflow-hidden">
                    <div className="bg-muted/20 p-6 border-b border-muted flex items-center gap-4">
                        <Avatar className="h-16 w-16 border border-muted ring-2 ring-background ring-offset-1 ring-offset-background shadow-sm">
                            <AvatarFallback className="bg-primary/10 text-primary font-bold text-xl">
                                {user.name.substring(0, 2).toUpperCase()}
                            </AvatarFallback>
                        </Avatar>
                        <div>
                            <h2 className="text-xl font-bold">{user.name}</h2>
                            <p className="text-sm text-muted-foreground">
                                {user.createdAt ? `Member since ${formatTimeAgo(user.createdAt)}` : "Circle member"}
                            </p>
                        </div>
                    </div>

                    <CardContent className="p-6 space-y-4">
                        <div className="flex items-center gap-3 text-sm">
                            <div className="h-9 w-9 rounded-lg bg-muted/50 flex items-center justify-center shrink-0">
                                <UserIcon className="h-4 w-4 text-muted-foreground" />
                            </div>
                            <div>
                                <div className="text-xs text-muted-foreground">Name</div>
                                <div className="font-medium">{user.name}</div>
                            </div>
                        </div>
                        <div className="flex items-center gap-3 text-sm">
                            <div className="h-9 w-9 rounded-lg bg-muted/50 flex items-center justify-center shrink-0">
                                <Mail className="h-4 w-4 text-muted-foreground" />
                            </div>
                            <div>
                                <div className="text-xs text-muted-foreground">Email</div>
                                <div className="font-medium">{user.email}</div>
                            </div>
                        </div>

                        <Button
                            variant="outline"
                            className="w-full gap-2 mt-4 text-destructive hover:text-destructive hover:bg-destructive/10"
                            onClick={logout}
                        >
                            <LogOut className="h-4 w-4" />
                            Log out
                        </Button>
                    </CardContent>
                </Card>
            )}
        </div>
    );
}
